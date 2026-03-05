--
-- PostgreSQL database dump
--

\restrict 6oOgFgM5Gr0ZbWB3jzVleb2fSxvSUGTAa03zoojPpfjToPANk9zTvhxRD4SLWce

-- Dumped from database version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)

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


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: inventories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inventories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    sku character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    price integer DEFAULT 0 NOT NULL,
    customer character varying(255),
    quantity integer DEFAULT 0 NOT NULL,
    reserved_quantity integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.inventories OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: stock_in; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_in (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    transaction_id character varying(255) NOT NULL,
    status character varying(50) DEFAULT 'created'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.stock_in OWNER TO postgres;

--
-- Name: stock_in_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_in_items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    stock_in_id uuid NOT NULL,
    inventory_id uuid NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.stock_in_items OWNER TO postgres;

--
-- Name: stock_out; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_out (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    transaction_id character varying(255) NOT NULL,
    status character varying(50) DEFAULT 'draft'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_stock_out_status CHECK (((status)::text = ANY ((ARRAY['DRAFT'::character varying, 'ALLOCATED'::character varying, 'IN_PROGRESS'::character varying, 'DONE'::character varying, 'CANCELLED'::character varying])::text[])))
);


ALTER TABLE public.stock_out OWNER TO postgres;

--
-- Name: stock_out_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_out_items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    stock_out_id uuid NOT NULL,
    inventory_id uuid NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.stock_out_items OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: inventories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventories (id, sku, name, price, customer, quantity, reserved_quantity, created_at, updated_at) FROM stdin;
2ea80192-c25e-4309-b463-93eec1ce2c31	SKU0005	Lenovo Yoga 3	12000000	\N	100	0	2026-03-05 00:48:20.592426	2026-03-05 03:36:10.147649
7cde423e-98ca-48ea-958a-ceb5595552c2	TEST-001	Test Product	100	\N	100	0	2026-03-04 20:47:38.360524	2026-03-05 03:36:25.174275
f6fed99c-0349-4eac-a1c0-f036ec7ef7b2	SKU0002	SMK-25 Mini - Black	500000	\N	130	0	2026-03-04 16:26:29	2026-03-05 03:19:10.231669
5585b14c-82d5-4a5e-93fa-ab0bc33dbca5	SKU0001	DS Orca 24bit Gen 4	480000	\N	115	0	2026-03-04 16:25:16	2026-03-05 03:19:10.233739
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
9	f
\.


--
-- Data for Name: stock_in; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock_in (id, transaction_id, status, created_at, updated_at) FROM stdin;
ee096a6f-1f29-4c04-a365-e1a9023f28df	STI-20260304193743	done	2026-03-04 19:37:43.279819	2026-03-04 20:05:54.181747
f3625888-534b-4a84-a137-4e07cc483349	STI-20260305021100	done	2026-03-05 02:11:00.039149	2026-03-05 02:14:27.348672
c0b4bfe0-87ab-43a6-83a2-a595dfd6a2d1	STI-20260305025724	done	2026-03-05 02:57:24.466517	2026-03-05 02:57:27.571311
a1ef199d-6598-4e84-a37f-39688e557cc3	STI-20260304191100	done	2026-03-04 19:11:00.522499	2026-03-05 03:19:08.260309
160fadfd-3cb1-44ab-ac6b-9d57a14f74cb	STI-20260304191116	done	2026-03-04 19:11:16.271354	2026-03-05 03:19:08.88477
1872912a-f9b3-4bfd-944a-179cecedc496	STI-20260304190939	done	2026-03-04 19:09:39.758376	2026-03-05 03:19:10.234824
\.


--
-- Data for Name: stock_in_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock_in_items (id, stock_in_id, inventory_id, quantity, created_at, updated_at) FROM stdin;
88ec5a1c-61a5-45f9-83c8-165e67f92170	1872912a-f9b3-4bfd-944a-179cecedc496	f6fed99c-0349-4eac-a1c0-f036ec7ef7b2	10	2026-03-04 19:09:39.761944	2026-03-04 19:09:39.761944
a64bdaa2-6aa3-4155-9e08-2b9f953acdec	1872912a-f9b3-4bfd-944a-179cecedc496	5585b14c-82d5-4a5e-93fa-ab0bc33dbca5	5	2026-03-04 19:09:39.761944	2026-03-04 19:09:39.761944
35ad2566-273f-4804-b6b1-0c9657f7024f	a1ef199d-6598-4e84-a37f-39688e557cc3	f6fed99c-0349-4eac-a1c0-f036ec7ef7b2	10	2026-03-04 19:11:00.526119	2026-03-04 19:11:00.526119
e894349a-ae62-4b3c-8f48-a3d38c164e83	a1ef199d-6598-4e84-a37f-39688e557cc3	5585b14c-82d5-4a5e-93fa-ab0bc33dbca5	5	2026-03-04 19:11:00.526119	2026-03-04 19:11:00.526119
d630c9a3-a027-4924-b023-44db3f72bca7	160fadfd-3cb1-44ab-ac6b-9d57a14f74cb	f6fed99c-0349-4eac-a1c0-f036ec7ef7b2	10	2026-03-04 19:11:16.273118	2026-03-04 19:11:16.273118
2d9f97e1-4ecd-4755-92a2-882d9ab189d2	160fadfd-3cb1-44ab-ac6b-9d57a14f74cb	5585b14c-82d5-4a5e-93fa-ab0bc33dbca5	5	2026-03-04 19:11:16.273118	2026-03-04 19:11:16.273118
0b80d964-c51d-4025-815e-a80e722e92aa	ee096a6f-1f29-4c04-a365-e1a9023f28df	f6fed99c-0349-4eac-a1c0-f036ec7ef7b2	10	2026-03-04 19:37:43.282433	2026-03-04 19:37:43.282433
1194f743-d85c-4106-b867-921c172738d0	ee096a6f-1f29-4c04-a365-e1a9023f28df	5585b14c-82d5-4a5e-93fa-ab0bc33dbca5	5	2026-03-04 19:37:43.282433	2026-03-04 19:37:43.282433
fca21b66-5717-4555-88a9-a783bea38761	f3625888-534b-4a84-a137-4e07cc483349	2ea80192-c25e-4309-b463-93eec1ce2c31	5	2026-03-05 02:11:00.04157	2026-03-05 02:11:00.04157
9d7e5b8a-d919-4c7f-acdc-a4f9c1e1da30	f3625888-534b-4a84-a137-4e07cc483349	5585b14c-82d5-4a5e-93fa-ab0bc33dbca5	10	2026-03-05 02:11:00.04157	2026-03-05 02:11:00.04157
2ccb8aad-16dd-448e-aecd-6e8f0b4272e8	c0b4bfe0-87ab-43a6-83a2-a595dfd6a2d1	2ea80192-c25e-4309-b463-93eec1ce2c31	200	2026-03-05 02:57:24.469322	2026-03-05 02:57:24.469322
\.


--
-- Data for Name: stock_out; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock_out (id, transaction_id, status, created_at, updated_at) FROM stdin;
87a6f477-c1fb-46c6-bcf9-cb23a6f60a23	STO-20260305031326	DONE	2026-03-05 03:13:26.400872	2026-03-05 03:15:17.531836
5c1be1a1-5e8f-4643-aa5b-5f97633d4cb4	STO-20260305031523	DONE	2026-03-05 03:15:23.45654	2026-03-05 03:18:08.578762
67a95b74-628d-464f-bff8-ba673916677a	STO-20260305031751	DONE	2026-03-05 03:17:51.170992	2026-03-05 03:18:09.042462
a016dd40-a986-4e85-8219-0910738c5346	STO-20260305031817	DONE	2026-03-05 03:18:17.790219	2026-03-05 03:18:43.907483
dbf13435-bcc9-4600-a87d-7a8d36ad6229	STO-20260305031826	DONE	2026-03-05 03:18:26.35969	2026-03-05 03:18:45.102876
6df1f141-7c40-4fb9-9202-efe634f15111	STO-20260305032339	IN_PROGRESS	2026-03-05 03:23:39.615787	2026-03-05 03:23:48.762821
\.


--
-- Data for Name: stock_out_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock_out_items (id, stock_out_id, inventory_id, quantity, created_at, updated_at) FROM stdin;
92a225eb-f714-439c-ab75-da1277ba2d8e	87a6f477-c1fb-46c6-bcf9-cb23a6f60a23	2ea80192-c25e-4309-b463-93eec1ce2c31	2	2026-03-05 03:13:26.402836	2026-03-05 03:13:26.402836
e1499215-d610-40a8-8ed1-d15c6ecc1910	5c1be1a1-5e8f-4643-aa5b-5f97633d4cb4	2ea80192-c25e-4309-b463-93eec1ce2c31	2	2026-03-05 03:15:23.459185	2026-03-05 03:15:23.459185
32f2d1f1-4e4a-4c69-9ecd-e5d5038fcecb	67a95b74-628d-464f-bff8-ba673916677a	7cde423e-98ca-48ea-958a-ceb5595552c2	2	2026-03-05 03:17:51.17368	2026-03-05 03:17:51.17368
16520cd1-344f-4cd9-94f8-ea3289552bec	a016dd40-a986-4e85-8219-0910738c5346	7cde423e-98ca-48ea-958a-ceb5595552c2	2	2026-03-05 03:18:17.792054	2026-03-05 03:18:17.792054
07bf2f24-3646-469b-a580-6309bf6cec6e	dbf13435-bcc9-4600-a87d-7a8d36ad6229	7cde423e-98ca-48ea-958a-ceb5595552c2	4	2026-03-05 03:18:26.361367	2026-03-05 03:18:26.361367
192ad7ec-5dc1-4ce2-93a8-42a8bfa84e3d	6df1f141-7c40-4fb9-9202-efe634f15111	2ea80192-c25e-4309-b463-93eec1ce2c31	2	2026-03-05 03:23:39.617654	2026-03-05 03:23:39.617654
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password, created_at) FROM stdin;
19e0b297-d039-4121-b88f-61522dd90dba	Irfan Ardiansyah	dev@mail.com	$2a$10$qs6IpFsJJWggosSXhVKpveWRlqEMFtEH8gJCS8h5WQmi4rNOd/Kra	2026-03-04 16:24:57.905935
7a0d2f0c-8bb8-4593-849c-969315298b79	Test User	test@example.com	$2a$10$dyeMMLt7K4OAZr.rkezF3eStAHSz4bW2Lni/juvIpNsDUhgb.nFvW	2026-03-04 20:46:42.750964
\.


--
-- Name: inventories inventories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories
    ADD CONSTRAINT inventories_pkey PRIMARY KEY (id);


--
-- Name: inventories inventories_sku_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories
    ADD CONSTRAINT inventories_sku_key UNIQUE (sku);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: stock_in_items stock_in_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_in_items
    ADD CONSTRAINT stock_in_items_pkey PRIMARY KEY (id);


--
-- Name: stock_in stock_in_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_in
    ADD CONSTRAINT stock_in_pkey PRIMARY KEY (id);


--
-- Name: stock_out_items stock_out_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_out_items
    ADD CONSTRAINT stock_out_items_pkey PRIMARY KEY (id);


--
-- Name: stock_out stock_out_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_out
    ADD CONSTRAINT stock_out_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_inventory_sku; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_inventory_sku ON public.inventories USING btree (sku);


--
-- Name: idx_stock_in_items_inventory_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_in_items_inventory_id ON public.stock_in_items USING btree (inventory_id);


--
-- Name: idx_stock_in_items_stock_in_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_in_items_stock_in_id ON public.stock_in_items USING btree (stock_in_id);


--
-- Name: idx_stock_out_items_inventory_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_out_items_inventory_id ON public.stock_out_items USING btree (inventory_id);


--
-- Name: idx_stock_out_items_stock_out_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_out_items_stock_out_id ON public.stock_out_items USING btree (stock_out_id);


--
-- Name: stock_in_items stock_in_items_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_in_items
    ADD CONSTRAINT stock_in_items_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventories(id) ON DELETE CASCADE;


--
-- Name: stock_in_items stock_in_items_stock_in_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_in_items
    ADD CONSTRAINT stock_in_items_stock_in_id_fkey FOREIGN KEY (stock_in_id) REFERENCES public.stock_in(id) ON DELETE CASCADE;


--
-- Name: stock_out_items stock_out_items_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_out_items
    ADD CONSTRAINT stock_out_items_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventories(id) ON DELETE CASCADE;


--
-- Name: stock_out_items stock_out_items_stock_out_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_out_items
    ADD CONSTRAINT stock_out_items_stock_out_id_fkey FOREIGN KEY (stock_out_id) REFERENCES public.stock_out(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict 6oOgFgM5Gr0ZbWB3jzVleb2fSxvSUGTAa03zoojPpfjToPANk9zTvhxRD4SLWce

