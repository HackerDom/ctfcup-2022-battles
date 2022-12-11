create schema if not exists secrets;

create table secrets.users
(
    name  text constraint users_pk primary key,
    token text
);

create table secrets.secrets
(
    id       uuid default gen_random_uuid()
        constraint secrets_pk
            primary key,
    owner    text
        constraint secrets_users_name_fk
            references secrets.users,
    secret text
);

INSERT INTO secrets.users (name, token) VALUES ('admin', 'VmKps69TRvAFKP6uxzFdFAtVWs4Rw');
INSERT INTO secrets.secrets (owner, secret) VALUES ('admin', 'Cup{y0u_haCk_mY_seRv1cE_a3_Wha3_c0s3}')