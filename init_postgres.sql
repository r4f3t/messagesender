-- Create the example_table
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    "to" VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



INSERT INTO messages ("to", content, sent) VALUES
('+905551111111', 'Test Message 1', FALSE),
('+905552222222', 'Test Message 2', FALSE),
('+905553333333', 'Test Message 3', FALSE),
('+905554444444', 'Test Message 4', FALSE),
('+905555555555', 'Test Message 5', FALSE);

