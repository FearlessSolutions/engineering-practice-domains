-- migrate:up

--
-- Table structure for table `greetings`
--
CREATE TABLE greetings
(
    `id` int auto_increment PRIMARY KEY NOT NULL,
    `greetingText` varchar(32) NOT NULL
);

INSERT INTO greetings(greetingText) VALUES
    ('Hi'),
    ('Hello'),
    ('Howdy-do');

-- migrate:down
DROP TABLE IF EXISTS greetings;
