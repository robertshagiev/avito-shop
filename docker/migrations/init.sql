CREATE TABLE IF NOT EXISTS public.users (
                            id BIGSERIAL PRIMARY KEY,
                            username VARCHAR(100) UNIQUE NOT NULL,
                            password TEXT NOT NULL,
                            created_at TIMESTAMP DEFAULT NOW(),
                            coins INT NOT NULL DEFAULT 1000 CHECK (coins >= 0)
);

CREATE TABLE IF NOT EXISTS public.transactions (
                            id BIGSERIAL PRIMARY KEY,
                            from_user_id BIGINT REFERENCES public.users(id) ON DELETE SET NULL,
                            to_user_id BIGINT REFERENCES public.users(id) ON DELETE SET NULL,
                            quantity INT NOT NULL CHECK (quantity > 0),
                            created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.merch (
                            id BIGSERIAL PRIMARY KEY,
                            name VARCHAR(100) UNIQUE NOT NULL,
                            price INT NOT NULL CHECK (price > 0)
);

CREATE TABLE IF NOT EXISTS public.inventory (
                            id SERIAL PRIMARY KEY,
                            user_id INT REFERENCES users(id) ON DELETE SET NULL,
                            merch_id INT REFERENCES merch(id) ON DELETE CASCADE,
                            quantity INT NOT NULL DEFAULT 1,
                            bought_on TIMESTAMP DEFAULT NOW()
);

INSERT INTO public.merch (name, price) VALUES
                            ('t-shirt', 80),
                            ('cup', 20),
                            ('book', 50),
                            ('pen', 10),
                            ('powerbank', 200),
                            ('hoody', 300),
                            ('umbrella', 200),
                            ('socks', 10),
                            ('wallet', 50),
                            ('pink-hoody', 500)
ON CONFLICT (name) DO NOTHING;