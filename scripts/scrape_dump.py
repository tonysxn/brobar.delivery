
import os
import requests
from bs4 import BeautifulSoup
import uuid
import re

# Categories to scrape with their IDs (generating new UUIDs for consistency if needed, but keeping static if desired)
# We will generate UUIDs on the fly but keep them consistent for the run.
CATEGORIES = [
    {"name": "Бургери", "slug": "burgers"},
    {"name": "Сети", "slug": "sets"},
    {"name": "Перші страви", "slug": "persi stravi"},
    {"name": "Соуси", "slug": "sauces"},
    {"name": "Напої", "slug": "drinks"},
    {"name": "Фрі та сир", "slug": "fries-and-cheese"},
    {"name": "Салати", "slug": "salad"},
    {"name": "Гарячі страви", "slug": "hot"},
    {"name": "Гарячі закуски", "slug": "snacks"},
    {"name": "Холодні закуски", "slug": "holodni"},
]

BASE_URL = "https://brobar.delivery"
UPLOADS_DIR = "uploads_backup"


# Ensure uploads directory exists
if not os.path.exists(UPLOADS_DIR):
    os.makedirs(UPLOADS_DIR)

def clean_price(price_str):
    # Remove ' ₴', spaces, replace ',' with '.'
    val = re.sub(r'[^\d\.]', '', price_str.replace(',', '.'))
    try:
        return float(val)
    except:
        return 0.0

def clean_weight(weight_str):
    # Remove ' кг', etc
    val = re.sub(r'[^\d\.]', '', weight_str.replace(',', '.'))
    try:
        return float(val)
    except:
        return 0.0

def download_image(url, filename):
    try:
        response = requests.get(url, stream=True)
        if response.status_code == 200:
            # Clean filename: remove the hash part if possible, or just use the slug name
            # URL: .../products/slug.hash.hash.jpg
            # filename arg passed in is the slug (e.g. big-bro-burger)
            
            ext = 'jpg'
            if url.endswith('.png'):
                ext = 'png'
            elif url.endswith('.jpeg'):
                ext = 'jpeg'
                
            clean_name = f"{filename}.{ext}"
            filepath = os.path.join(UPLOADS_DIR, clean_name)
            
            with open(filepath, 'wb') as f:
                for chunk in response.iter_content(1024):
                    f.write(chunk)
            
            return clean_name

    except Exception as e:
        print(f"Failed to download {url}: {e}")
    return None

def main():
    print("Starting scrape...")
    
    products_sql = []
    categories_sql = []
    
    # Header for the SQL dump
    header = """--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

SET default_tablespace = '';

SET default_table_access_method = heap;

DROP TABLE IF EXISTS public.products CASCADE;
DROP TABLE IF EXISTS public.product_variations CASCADE;
DROP TABLE IF EXISTS public.product_variation_groups CASCADE;
DROP TABLE IF EXISTS public.categories CASCADE;
DROP TABLE IF EXISTS public.schema_migrations CASCADE;

CREATE TABLE public.categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    icon character varying(255),
    sort integer DEFAULT 0
);

ALTER TABLE public.categories OWNER TO sanin;

CREATE TABLE public.product_variation_groups (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    product_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    external_id character varying(100) NOT NULL,
    default_value integer,
    show boolean DEFAULT true,
    required boolean DEFAULT false
);

ALTER TABLE public.product_variation_groups OWNER TO sanin;

CREATE TABLE public.product_variations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    group_id uuid NOT NULL,
    external_id character varying(100) NOT NULL,
    default_value integer,
    show boolean DEFAULT true,
    name character varying(255) NOT NULL
);

ALTER TABLE public.product_variations OWNER TO sanin;

CREATE TABLE public.products (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    external_id character varying(100) NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    weight numeric(10,3) DEFAULT 0,
    category_id uuid,
    sort integer DEFAULT 0,
    hidden boolean DEFAULT false,
    alcohol boolean DEFAULT false,
    sold boolean DEFAULT false,
    image character varying(255) NOT NULL,
    CONSTRAINT products_price_check CHECK ((price >= (0)::numeric))
);

ALTER TABLE public.products OWNER TO sanin;

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

ALTER TABLE public.schema_migrations OWNER TO sanin;

COPY public.categories (id, name, slug, icon, sort) FROM stdin;
"""
    
    product_rows = []
    category_rows = []
    
    for idx, cat in enumerate(CATEGORIES):
        cat_id = str(uuid.uuid4())
        cat_name = cat['name']
        cat_slug = cat['slug']
        # category icon is not easily scraped from the card, but we can try or leave blank
        # The card has <i class="fal fa-hamburger"></i> etc.
        # We'll just leave it null for now or hardcode if essential.
        # Use empty string for icon to avoid NULL scan error in Go struct (Icon string)
        category_rows.append(f"{cat_id}\t{cat_name}\t{cat_slug}\t\t{idx*10}")
        
        url = f"{BASE_URL}/category/{cat['slug']}"
        print(f"Scraping category: {cat_name} ({url})")
        
        try:
            resp = requests.get(url)
            if resp.status_code != 200:
                print(f"Failed to fetch {url}")
                continue
                
            soup = BeautifulSoup(resp.content, 'html.parser')
            product_cards = soup.select('.card-product')
            
            for prod_card in product_cards:
                try:
                    # Name
                    name_tag = prod_card.select_one('.card-product-name p')
                    name = name_tag.get_text(strip=True) if name_tag else "Unknown"

                    # Slug (derive from name or link)
                    link_tag = prod_card.find_parent('a')
                    href = link_tag['href'] if link_tag else ""
                    # href: https://brobar.delivery/product/big-bro-burger
                    slug = href.split('/')[-1] if href else uuid.uuid4().hex

                    # Image
                    img_tag = prod_card.select_one('div.card-product-img img')
                    img_url = img_tag['src'] if img_tag else ""
                    if img_url:
                        # Ensure it's absolute
                        if not img_url.startswith('http'):
                            img_url = BASE_URL + img_url
                        
                        # Pass slug as the desired filename base
                        img_filename = download_image(img_url, slug)

                    else:
                        img_filename = ""
                    

                    # Price and Weight
                    price = 0.0
                    weight = 0.0
                    info_bg = prod_card.select_one('.card-product-price-bg')
                    if info_bg:
                        # Text usually: "0.600 кг | 470 ₴"
                        text_content = info_bg.get_text(strip=True, separator="|")
                        parts = text_content.split('|')
                        for p in parts:
                            p = p.strip()
                            if '₴' in p:
                                price = clean_price(p)
                            if 'кг' in p or 'л' in p: # handle liters if any
                                weight = clean_weight(p)
                    
                    # Description
                    desc_tag = prod_card.select_one('.card-product-info-hover-text p')
                    description = desc_tag.get_text(strip=True) if desc_tag else ""
                    description = desc_tag.get_text(strip=True) if desc_tag else ""
                    # Escape tabs and newlines for COPY format. Remove \r completely.
                    description = description.replace('\r', '').replace('\n', '\\n').replace('\t', ' ')



                    
                    prod_id = str(uuid.uuid4())
                    external_id = str(uuid.uuid4())[:8] # Fake external ID
                    
                    # Row format: id, external_id, name, slug, description, price, weight, category_id, sort, hidden, alcohol, sold, image
                    row = f"{prod_id}\t{external_id}\t{name}\t{slug}\t{description}\t{price}\t{weight}\t{cat_id}\t0\tf\tf\tf\t{img_filename}"
                    product_rows.append(row)
                    
                except Exception as e:
                    print(f"Error processing product in {cat_name}: {e}")
                    
        except Exception as e:
            print(f"Error scraping {url}: {e}")

    # Assemble the file content
    content = header
    content += "\n".join(category_rows)
    content += "\n\\.\n\n"
    
    content += """--
-- Data for Name: product_variation_groups; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.product_variation_groups (id, product_id, name, external_id, default_value, show, required) FROM stdin;
\\.


--
-- Data for Name: product_variations; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.product_variations (id, group_id, external_id, default_value, show, name) FROM stdin;
\\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.products (id, external_id, name, slug, description, price, weight, category_id, sort, hidden, alcohol, sold, image) FROM stdin;
"""
    content += "\n".join(product_rows)
    content += "\n\\.\n\n"
    
    content += """--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.schema_migrations (version, dirty) FROM stdin;
3	f
\\.


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: categories categories_slug_key; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);


--
-- Name: product_variation_groups product_variation_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.product_variation_groups
    ADD CONSTRAINT product_variation_groups_pkey PRIMARY KEY (id);


--
-- Name: product_variations product_variations_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.product_variations
    ADD CONSTRAINT product_variations_pkey PRIMARY KEY (id);


--
-- Name: products products_external_id_key; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_external_id_key UNIQUE (external_id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: products products_slug_key; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_slug_key UNIQUE (slug);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: product_variation_groups uq_product_variation_groups_product_external_id; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.product_variation_groups
    ADD CONSTRAINT uq_product_variation_groups_product_external_id UNIQUE (product_id, external_id);


--
-- Name: product_variations uq_product_variations_group_external_id; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.product_variations
    ADD CONSTRAINT uq_product_variations_group_external_id UNIQUE (group_id, external_id);


--
-- Name: idx_product_variation_groups_product_id; Type: INDEX; Schema: public; Owner: sanin
--

CREATE INDEX idx_product_variation_groups_product_id ON public.product_variation_groups USING btree (product_id);


--
-- Name: idx_product_variations_group_id; Type: INDEX; Schema: public; Owner: sanin
--

CREATE INDEX idx_product_variations_group_id ON public.product_variations USING btree (group_id);


--
-- Name: idx_products_category; Type: INDEX; Schema: public; Owner: sanin
--

CREATE INDEX idx_products_category ON public.products USING btree (category_id);


--
-- Name: idx_products_slug; Type: INDEX; Schema: public; Owner: sanin
--

CREATE INDEX idx_products_slug ON public.products USING btree (slug);


--
-- Name: product_variation_groups product_variation_groups_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.product_variation_groups
    ADD CONSTRAINT product_variation_groups_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: product_variations product_variations_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.product_variations
    ADD CONSTRAINT product_variations_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.product_variation_groups(id) ON DELETE CASCADE;


--
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--
"""

    with open('backup_product_db.sql', 'w') as f:
        f.write(content)
        print("Created backup_product_db.sql")

if __name__ == "__main__":
    main()
