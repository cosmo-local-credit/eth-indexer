-- Nullable without a default: rows indexed from the legacy Swap event carry no
-- settlement detail, and a zero would be indistinguishable from a real zero fee.
ALTER TABLE pool_swap
ADD COLUMN quoted_out_value NUMERIC,
ADD COLUMN nominal_out_value NUMERIC,
ADD COLUMN protocol_fee NUMERIC;
