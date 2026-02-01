--
-- PostgreSQL database dump
--

\restrict 0xElogK2Yc8zhKwOryKGmtdcVioFKSpBCba89SrLoOV2mPzs082EQg8qqpPvpw0

-- Dumped from database version 14.20 (Debian 14.20-1.pgdg13+1)
-- Dumped by pg_dump version 14.20 (Debian 14.20-1.pgdg13+1)

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

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: user_role; Type: TYPE; Schema: public; Owner: sanin
--

CREATE TYPE public.user_role AS ENUM (
    'admin',
    'user',
    'moderator'
);


ALTER TYPE public.user_role OWNER TO sanin;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: sanin
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO sanin;

--
-- Name: users; Type: TABLE; Schema: public; Owner: sanin
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    role_id public.user_role DEFAULT 'user'::public.user_role NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    name text NOT NULL,
    address text,
    address_coords text,
    phone text,
    promo_card text
);


ALTER TABLE public.users OWNER TO sanin;

--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.schema_migrations (version, dirty) FROM stdin;
1	f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.users (id, role_id, email, password, name, address, address_coords, phone, promo_card) FROM stdin;
8b5d58d2-4511-4c37-bfd1-ed2a46469a96	admin	sanin.tony.dev@gmail.com	$2a$10$noQ.DLtzEJt/FvFUDueIYOxGrmJIMZ.WkJRL3iXrC/HcQibV5mYs6	Admin	\N	\N	\N	\N
\.


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: sanin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict 0xElogK2Yc8zhKwOryKGmtdcVioFKSpBCba89SrLoOV2mPzs082EQg8qqpPvpw0

