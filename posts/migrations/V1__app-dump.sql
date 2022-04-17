--
-- PostgreSQL database dump
--

-- Dumped from database version 12.9 (Debian 12.9-1.pgdg110+1)
-- Dumped by pg_dump version 12.9 (Ubuntu 12.9-0ubuntu0.20.04.1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: posts; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.posts (
    id integer NOT NULL,
    name character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    deleted boolean
);


ALTER TABLE public.posts OWNER TO root;

--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.posts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_id_seq OWNER TO root;

--
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;


--
-- Name: posts id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);


--
-- Data for Name: posts; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.posts (id, name, created_at, updated_at, deleted) FROM stdin;
10	asa	2022-04-15 06:44:26.922568	2022-04-15 06:44:26.922568	f
9	asa	2022-04-15 06:44:25.717996	2022-04-15 06:44:25.717996	f
11	asa	2022-04-15 06:44:28.004724	2022-04-15 06:44:28.004724	f
\.


--
-- Name: posts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.posts_id_seq', 11, true);


--
-- Name: posts posts_pk; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

