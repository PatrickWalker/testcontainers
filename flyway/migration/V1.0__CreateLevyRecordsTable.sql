CREATE TABLE users
(
    id                               SERIAL,
    firstname                             VARCHAR,
    middlename                        VARCHAR,
    lastname                VARCHAR,
    title                  VARCHAR,
    dob                              VARCHAR,
    update_time                      TIMESTAMP WITH TIME ZONE,
    insert_time                      TIMESTAMP WITH TIME ZONE DEFAULT now(),
    PRIMARY KEY(id)
);

INSERT INTO users("firstname","middlename","lastname","title","dob","update_time")
VALUES('don','pablo','jr','mr','01/01/1960',now());