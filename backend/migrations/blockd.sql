create table if not exists users (
        id uuid primary key ,
        name varchar(250),
        email varchar(200),
        phone varchar(16),
        tg varchar(200),
        public_key bytea not null unique,
        mnemonic varchar(500) not null,
        seed bytea not null unique,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp,
        activated_at  timestamp default null
);

create index if not exists index_users_seed
        on users (seed); 

create index if not exists index_users_public_key
        on users (public_key); 

create index if not exists index_users_name
        on users using hash (name); 

create index if not exists index_users_seed
        on users using hash (seed); 

create table if not exists access_tokens (
        user_id uuid not null references users(id),
        token varchar(350) not null, 
        token_expired_at timestamp, 
        refresh_token varchar(350) not null, 
        refresh_token_expired_at timestamp, 
        created_at timestamp default current_timestamp,
        remote_addr varchar(100)
);

create index if not exists index_access_tokens_token_refresh_token
        on access_tokens (token, refresh_token); 

create index if not exists index_access_tokens_token_refresh_token_exp
        on access_tokens (token, refresh_token, token_expired_at, refresh_token_expired_at); 

create table if not exists organizations (
        id uuid primary key unique, 
        name varchar(300) default 'My Organization' not null, 
        address varchar(750) not null, 
        wallet_seed bytea not null,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
);

create index if not exists index_organizations_id
        on organizations (id); 

create table employees (
        id uuid primary key, 
        user_id uuid, 
        organization_id uuid not null references organizations(id),
        wallet_address text not null, 
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
);

create index if not exists index_employees_id_organization_id
        on employees (id, organization_id); 

create index if not exists index_user_id_organization_id
        on employees (user_id, organization_id); 

create table organizations_users (
        organization_id uuid not null references organizations(id), 
        user_id uuid not null references users(id), 
        employee_id uuid default null,
        position varchar(300),
        added_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp,
        deleted_at timestamp default null,
        is_admin bool default false,
        is_owner bool default false,
        primary key(organization_id, user_id)
);

create index if not exists index_organizations_users_organization_id_user_id_is_admin
        on organizations_users (organization_id, user_id, is_admin); 

create index if not exists index_organizations_users_organization_id_user_id
        on organizations_users (organization_id, user_id); 

create index if not exists index_organizations_users_organization_id_employee_id
        on organizations_users (organization_id, employee_id); 


create table if not exists transactions (
        id uuid primary key,
        description text default 'New Transaction', 
        organization_id uuid not null, 
        created_by uuid  not null, 
        amount decimal default 0,

        to_addr bytea not null,
        tx_index bytea default null,

        max_fee_allowed decimal default 0, 
        deadline timestamp default null,
        status int default 0,

        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp,

        confirmed_at timestamp default null,
        cancelled_at timestamp default null,

        commited_at timestamp default null
);

create index if not exists index_transactions_id_organization_id
        on transactions (organization_id); 

create index if not exists index_transactions_id_organization_id_created_by
        on transactions (organization_id, created_by); 

create index if not exists index_transactions_organization_id_deadline
        on transactions (organization_id, deadline); 

create table transactions_confirmations (
        tx_id uuid not null, 
        user_id uuid not null,
        organization_id uuid not null, 
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp,
        confirmed bool
);

create index if not exists index_transactions_confirmations_tx_id_user_id_organization_id
        on transactions_confirmations (tx_id, user_id, organization_id);

create table contracts (
        id uuid primary key, 
        title varchar(250) default 'New Contract', 
        description text not null, 

        address bytea not null, 

        payload bytea not null,

        created_by uuid not null references users(id), 
        organization_id uuid not null references organizations(id), 

        status smallint default 0,
        tx_index bytea default null,

        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
);

create table multisigs (
        id uuid primary key, 
        title varchar(350) default 'New Multi-Sig'
);

create table multisig_owners (
        multisig_id uuid references multisigs(id), 
        owner_id uuid references users(id), 
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp,
        primary key (multisig_id, owner_id)
);

create table multisig_confirmations (
        multisig_id uuid references multisigs(id), 
        owner_id uuid references users(id), 
        created_at timestamp default current_timestamp,
        primary key (multisig_id, owner_id)
);