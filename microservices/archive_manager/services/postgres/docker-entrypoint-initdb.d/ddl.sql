CREATE TABLE record (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "text" TEXT
);

ALTER TABLE record
ADD CONSTRAINT record_unique_name CHECK (LENGTH("text")>=10 AND LENGTH("text") <= 500);

ALTER TABLE record
ALTER COLUMN "text" SET NOT NULL;