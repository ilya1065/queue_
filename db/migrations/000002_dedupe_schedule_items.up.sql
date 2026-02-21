DELETE FROM schedule_items
WHERE id NOT IN (
    SELECT MIN(id)
    FROM schedule_items
    GROUP BY name, start_date, end_date
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_schedule_items_name_start_end
    ON schedule_items(name, start_date, end_date);
