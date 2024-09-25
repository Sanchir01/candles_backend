--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.4 (Debian 16.4-1.pgdg120+1)

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
-- Name: update_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_timestamp() OWNER TO postgres;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: candles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.candles (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title text NOT NULL,
    slug text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    category_id uuid NOT NULL,
    images text[] NOT NULL,
    price numeric DEFAULT 0,
    color_id uuid
);


ALTER TABLE public.candles OWNER TO postgres;

--
-- Name: category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.category (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.category OWNER TO postgres;

--
-- Name: color; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.color (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.color OWNER TO postgres;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE public.goose_db_version OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.goose_db_version_id_seq OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    role character varying(50) DEFAULT 'user'::character varying,
    phone character varying(11) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Data for Name: candles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.candles (id, title, slug, created_at, updated_at, version, category_id, images, price, color_id) FROM stdin;
afe7f9e5-6a7c-4593-bca5-c2ce0e296a16	mix candles	mix-candles	2024-09-09 14:32:01.14485	2024-09-09 14:32:01.14485	1	801896ce-f22b-4006-8ddb-752375cebe08	{https://unsplash.com/photos/a-candle-sitting-on-top-of-a-book-on-a-table-mixRqgmyXAc,https://unsplash.com/photos/white-candle-on-black-holder-0vdLxhrb1Qw}	0	d7862a18-e4dd-41dd-a15c-8ee67b30a0c6
6c84bde6-6887-4ec6-94ec-f7d0e02a4225	black candle	black-candle	2024-09-10 06:54:58.693803	2024-09-10 06:54:58.693803	1	801896ce-f22b-4006-8ddb-752375cebe08	{https://unsplash.com/photos/a-candle-sitting-on-top-of-a-book-on-a-table-mixRqgmyXAc,https://unsplash.com/photos/white-candle-on-black-holder-0vdLxhrb1Qw}	3000	d7862a18-e4dd-41dd-a15c-8ee67b30a0c6
\.


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.category (id, title, slug, created_at, updated_at, version) FROM stdin;
21ebdaf3-a217-40f4-92ee-d9e235e9cebb	sanchir	sanchir	2024-09-09 14:31:39.458115	2024-09-09 14:31:39.458115	1
801896ce-f22b-4006-8ddb-752375cebe08	candle	candle	2024-09-09 14:31:48.036379	2024-09-09 14:31:48.036379	1
ffb11fba-f4b1-4bec-9b75-e68bde3a84ac	test	test	2024-09-10 08:09:51.650525	2024-09-10 08:09:51.650525	1
d97cc81a-83bc-49f8-84e0-a811b8ebdd7c	candles	candles	2024-09-12 08:49:09.628346	2024-09-12 08:49:09.628346	1
310362b6-5506-4be6-87df-9e4a57a819c7	abc	abc	2024-09-23 12:22:22.25803	2024-09-23 12:22:22.25803	1
\.


--
-- Data for Name: color; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.color (id, title, slug, created_at, updated_at, version) FROM stdin;
d7862a18-e4dd-41dd-a15c-8ee67b30a0c6	red	red	2024-09-13 10:00:41.304348	2024-09-13 10:00:41.304348	1
71632435-1dc5-4f3b-b0bc-08fb4e8d7c61	black	black	2024-09-13 10:00:47.788017	2024-09-13 10:00:47.788017	1
\.


--
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
1	0	t	2024-09-02 13:36:30.977746
2	20240830204746	t	2024-09-02 13:53:15.437274
3	20240830212655	t	2024-09-05 10:00:21.749389
4	20240902184055	t	2024-09-05 10:14:28.78314
5	20240905095214	t	2024-09-05 10:14:28.826486
6	20240905120352	t	2024-09-05 12:11:21.513952
7	20240906101526	t	2024-09-06 10:19:20.587095
8	20240910133503	t	2024-09-10 14:01:02.781906
9	20240910191634	t	2024-09-11 06:56:58.658269
10	20240910192341	t	2024-09-11 06:56:58.732366
11	20240910193833	t	2024-09-11 06:56:58.737842
12	20240910200526	t	2024-09-11 06:56:58.749571
13	20240911073646	t	2024-09-11 07:40:39.284806
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, title, slug, created_at, updated_at, version, role, phone) FROM stdin;
\.


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 13, true);


--
-- Name: candles candles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.candles
    ADD CONSTRAINT candles_pkey PRIMARY KEY (id);


--
-- Name: candles candles_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.candles
    ADD CONSTRAINT candles_slug_key UNIQUE (slug);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (id);


--
-- Name: category category_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_slug_key UNIQUE (slug);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: color size_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.color
    ADD CONSTRAINT size_pkey PRIMARY KEY (id);


--
-- Name: color size_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.color
    ADD CONSTRAINT size_slug_key UNIQUE (slug);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_slug_key UNIQUE (slug);


--
-- Name: color set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.color FOR EACH ROW EXECUTE FUNCTION public.update_timestamp();


--
-- Name: users set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_timestamp();


--
-- Name: category update_category_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_category_updated_at BEFORE UPDATE ON public.category FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: candles candles_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.candles
    ADD CONSTRAINT candles_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(id);


--
-- PostgreSQL database dump complete
--

