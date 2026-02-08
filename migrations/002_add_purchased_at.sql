-- Migration: 002_add_purchased_at
-- Description: Add purchased_at date field to assets table

ALTER TABLE assets
ADD COLUMN purchased_at DATE NULL AFTER notes;
