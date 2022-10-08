/* Setup for use with PostgreSQL */

/* STEP 1 + 2 */ 
create table if not exists user_info (
  id int primary key generated always as identity,
  first_name text,
  last_name text,
  city text,
  zip_code text
);
/*
  Setup table to auto delete when user is deleted,
  set a default change_date of now, and force the caller
  to set the password as currently_active if necessary,
  defaulting to false.
*/
create table if not exists user_password (
  id int primary key generated always as identity,
  user_info_id int references user_info on delete cascade,
  password_hash text not null,
  change_date timestamp default current_timestamp,
  currently_active boolean default false
);

/* Insert some base test data */
insert into user_info (first_name, last_name, city, zip_code) values 
  ('first', 'last', 'Boston', '12345'),
  ('Mary Anne', 'Test', 'New York', '01210'),
  ('Test', 'McTest', 'San Francisco', '00001'),
  ('UsesFull', 'ZipCode', 'San Francisco', '00001-1234');
insert into user_password (user_info_id, password_hash, currently_active)
  select
    id,
    /* not a password, but sufficient for testing */
    sha256((id || first_name || last_name || city)::bytea),
    TRUE
  from
    user_info;

/* Confirms deletion integrity */
select count(*) from user_password;
delete from user_info where first_name = 'Test';
select count(*) from user_password;

/* QUESTION 3 */
/* Selects all active passwords (or well, password hashes) */
select * from user_password where currently_active;

/* QUESTION 4, 5, and 6 */
/* Using the data above, show example of transaction to insert new password and update existing */
begin;
  /* First update all passwords for this user to be inactive */
  update
    user_password
  set
    currently_active = false
  where
    user_password.user_info_id in (
      select id from user_info where user_info.first_name = 'first' 
    )
    and user_password.currently_active;

  /* Then insert the new password */
  insert into user_password (user_info_id, password_hash, currently_active)
    select
      user_info.id,
      /* not a password, but sufficient for testing */
      sha256((user_info.id || user_info.first_name || user_info.last_name || user_info.city)::bytea),
      TRUE
    from
      user_info
    where
      user_info.first_name = 'first';
commit;

/* To confirm the above, confirm the first is incrementing, and the second is not */
select count(*) currently_inactive
from user_password
where
  user_password.user_info_id in (
    select id from user_info where user_info.first_name = 'first' 
  )
  and not user_password.currently_active;
select count(*)
from user_password
where
  user_password.user_info_id in (
    select id from user_info where user_info.first_name = 'first' 
  )
  and user_password.currently_active;