DELETE FROM short_link WHERE id NOT IN (
    SELECT MIN(id)
    FROM short_link
    GROUP BY url
);

ALTER TABLE short_link ADD CONSTRAINT unique_url UNIQUE (url);