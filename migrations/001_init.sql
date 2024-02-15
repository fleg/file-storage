create table files(
  id uuid primary key default gen_random_uuid() not null,
  uploaded_at timestamp default now() not null,
  size integer not null,
  -- https://datatracker.ietf.org/doc/html/rfc6838#section-4.2
  -- 127 + 1 + 127 == 255
  mime varchar(255) not null,
  -- usually it is 255 bytes across popular file systems
  -- https://en.wikipedia.org/wiki/Comparison_of_file_systems#Limits
  name varchar(255) not null
);

---- create above / drop below ----

drop table files;
