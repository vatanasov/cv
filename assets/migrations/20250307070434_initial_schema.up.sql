create table candidates
(
    id          integer primary key,
    given_name  text     not null,
    family_name text     not null,
    created_at  datetime not null default current_timestamp,
    updated_at  datetime not null default current_timestamp
);

create table communications
(
    id           integer primary key,
    candidate_id integer  not null,
    channel_code text     not null,
    uri          text     not null,
    created_at   datetime not null default current_timestamp,
    updated_at   datetime not null default current_timestamp,
    foreign key (candidate_id)
        references candidates (id)
        ON DELETE cascade
        ON UPDATE no action
);

create table employment_histories
(
    id                integer primary key,
    organization_name text     not null,
    position          text     not null,
    from_date         date     not null,
    to_date           date     not null,
    current           boolean  not null,
    text_description  text     not null,
    candidate_id      integer  not null,
    created_at        datetime not null default current_timestamp,
    updated_at        datetime not null default current_timestamp,
    foreign key (candidate_id)
        references candidates (id)
        ON DELETE cascade
        ON UPDATE no action
);

create table employment_history_roles
(
    id                    integer primary key,
    employment_history_id integer  not null,
    project               text     not null,
    role                  text     not null,
    created_at            datetime not null default current_timestamp,
    updated_at            datetime not null default current_timestamp,
    foreign key (employment_history_id)
        references employment_histories (id)
        ON DELETE cascade
        ON UPDATE no action
);

create table technologies
(
    name            text primary key,
    prettified_name text     not null,
    order_priority  integer  not null,
    created_at      datetime not null default current_timestamp,
    updated_at      datetime not null default current_timestamp,
    UNIQUE (name COLLATE NOCASE)
);

create table employment_history_roles_technologies
(
    employment_history_role_id integer  not null,
    technology_name            text     not null,
    created_at                 datetime not null default current_timestamp,
    updated_at                 datetime not null default current_timestamp,
    primary key (employment_history_role_id, technology_name),
    foreign key (employment_history_role_id)
        references employment_history_roles (id)
        ON DELETE cascade
        ON UPDATE no action,
    foreign key (technology_name)
        references technologies (name)
        ON DELETE cascade
        ON UPDATE no action
);

create table education_histories
(
    id                integer primary key,
    organization_name text     not null,
    degree            text     not null,
    candidate_id      integer  not null,
    from_date         date     not null,
    to_date           date     not null,
    created_at        datetime not null default current_timestamp,
    updated_at        datetime not null default current_timestamp,
    foreign key (candidate_id)
        references candidates (id)
        ON DELETE cascade
        ON UPDATE no action
);