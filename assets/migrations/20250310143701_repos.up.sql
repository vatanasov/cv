create table repos
(
    id              integer primary key,
    candidate_id    integer  not null,
    html_url        text     not null,
    name            text     not null,
    technology_name text     not null,
    created_at      datetime not null default current_timestamp,
    updated_at      datetime not null default current_timestamp,
    foreign key (candidate_id)
        references candidates (id)
        ON DELETE cascade
        ON UPDATE no action,
    foreign key (technology_name)
        references technologies (name)
        ON DELETE no action
        ON UPDATE no action,
    UNIQUE (html_url COLLATE NOCASE)
);
