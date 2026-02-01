--
-- PostgreSQL database dump
--

\restrict fdWLD7pqMRumDVfmH3D68B60S48KJUL4gAIBKCxVtX6sEYjsg0qysD4nZfobPa4

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


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: sanin
--

CREATE TABLE public.categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    icon character varying(255),
    sort integer DEFAULT 0
);


ALTER TABLE public.categories OWNER TO sanin;

--
-- Name: product_variation_groups; Type: TABLE; Schema: public; Owner: sanin
--

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

--
-- Name: product_variations; Type: TABLE; Schema: public; Owner: sanin
--

CREATE TABLE public.product_variations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    group_id uuid NOT NULL,
    external_id character varying(100) NOT NULL,
    default_value integer,
    show boolean DEFAULT true,
    name character varying(255) NOT NULL
);


ALTER TABLE public.product_variations OWNER TO sanin;

--
-- Name: products; Type: TABLE; Schema: public; Owner: sanin
--

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

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: sanin
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO sanin;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.categories (id, name, slug, icon, sort) FROM stdin;
f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	Бургери	burgers		0
524da3bb-bb79-4c53-ac60-a21dac89e202	Сети	sets		10
b400c7ab-817b-47c0-a72b-a66597d440e2	Перші страви	persi stravi		20
b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	Соуси	sauces		30
525976cc-90e9-497d-9be3-2ce001e718bf	Напої	drinks		40
aa5ace0d-00cb-4f37-a358-821d3afef2a1	Фрі та сир	fries-and-cheese		50
b5241ce2-58f7-42e3-a053-19105384ddb6	Салати	salad		60
2ed614c9-5d7a-4d45-ac4f-78589dd79848	Гарячі страви	hot		70
05362a6d-48e9-443b-a8c6-d405e7111a7c	Гарячі закуски	snacks		80
e3dabbf0-0e50-437d-bdd0-36e5b420082b	Холодні закуски	holodni		90
\.


--
-- Data for Name: product_variation_groups; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.product_variation_groups (id, product_id, name, external_id, default_value, show, required) FROM stdin;
\.


--
-- Data for Name: product_variations; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.product_variations (id, group_id, external_id, default_value, show, name) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.products (id, external_id, name, slug, description, price, weight, category_id, sort, hidden, alcohol, sold, image) FROM stdin;
7137b00d-7f7f-441e-b3b6-c3567e291466	e66920d7	BBQ	bbq	BBQ	30.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	bbq.webp
9f452ece-0080-4169-a2dd-2070c0a96d93	1f613dc8	BIG BRO burger	big-bro-burger	Подвійна котлета з телятини, подвійний cир чеддер, цибулеві кільця, солоний огірок, помідор, копчений BBQ, ромен	470.00	0.600	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	big-bro-burger.webp
a55f9b4f-159e-4805-80d7-59fb8397b8d0	6184fec7	Chicken burger	chicken-burger	Куряча котлета, бекон, помідор,  йогуртовий соус, сир чеддер, ромен	340.00	0.380	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	chicken-burger.webp
c83f90d7-4ba8-4664-9141-9fe3a1120555	c184ec4b	Clausthaler Б/А	clausthaler-ba	Clausthaler Б/А, Німеччина	100.00	0.500	525976cc-90e9-497d-9be3-2ce001e718bf	0	f	f	f	clausthaler-ba.webp
1be80dba-5a1a-493f-a9af-a27fbab3a897	ffc49c31	Фрі з беконовим мармеладом	fri-z-bekonovim-marmeladom	Картопля фрі, сирний соус, беконовий мармелад, кріспі цибуля, трюфельний майо, пармезан	235.00	0.270	aa5ace0d-00cb-4f37-a358-821d3afef2a1	0	f	f	f	fri-z-bekonovim-marmeladom.webp
b574680f-e99c-47ee-a17d-5b25c772cc27	5963eedd	Гірчиця	gircicia	Гірчиця	15.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	gircicia.webp
b81f1f1c-a301-4374-8989-2b1e01ef2b9d	4aac20e6	Картопля фрі	kartoplia-fri	Соус BBQ	120.00	0.170	aa5ace0d-00cb-4f37-a358-821d3afef2a1	0	f	f	f	kartoplia-fri.webp
1769f69d-503c-4f05-a551-f6249ab72965	5d2a882b	Картопляні діпи	kartopliani-dipi	Сирно-часниковий соус	140.00	0.170	aa5ace0d-00cb-4f37-a358-821d3afef2a1	0	f	f	f	kartopliani-dipi.webp
fcacd6b5-7645-4548-9570-4f5b80a3f1ff	f857fe15	Кетчуп	ketcup	Кетчуп	30.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	ketcup.webp
88be20f9-b26c-4475-98d3-5c8b9857654d	cc2a9e14	King Mushrooms burger	king-mushrooms-burger	Смеш котлета з телятини, два трюфельних  соуси, моцарела, солоний огірок, гриби ерінги	370.00	0.340	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	king-mushrooms-burger.webp
7e719368-62f7-4c08-9470-40ad06435a1d	e964e30e	Кока-кола	koka-kola	Кока-кола	60.00	0.500	525976cc-90e9-497d-9be3-2ce001e718bf	0	f	f	f	koka-kola.webp
4dc1a28a-f9b3-4d6a-8dde-1d71a4d00d43	0ae0ec40	Копчений BBQ	kopcenii-bbq	Копчений BBQ	40.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	kopcenii-bbq.webp
2e1d87ec-e20b-4884-9874-1cf2621eb6b8	0077fce6	Копчений майонез	kopcenii-maionez	Копчений майонез	30.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	kopcenii-maionez.webp
5b761ee0-db47-49af-8d57-8e64519576cf	96efe357	Майонез	maionez	Майонез	20.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	maionez.webp
8de522d5-2b4b-454d-ad90-bfad5f57e8f9	f8850af1	Meat burger	meat-burger	Котлета з телятини, свіжа цибуля, солоний огірок, помідор, сир чеддер, копчений майонез, ромен	340.00	0.360	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	meat-burger.webp
972f61cd-1aea-42ad-b8d3-c6853d284127	028a361a	Медово-гірчичний	medovo-gircicnii	Медово-гірчичний	30.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	medovo-gircicnii.webp
eb661642-c8cf-4d2c-9b00-a374308b7af6	aea18ff0	Моршинська газ	morsinska-gaz	Газ	40.00	0.500	525976cc-90e9-497d-9be3-2ce001e718bf	0	f	f	f	morsinska-gaz.webp
5e4139b0-541d-4bd0-97f9-1d31d231404c	3dc12c05	Моршинська негаз	morsinska-negaz	Негаз	40.00	0.500	525976cc-90e9-497d-9be3-2ce001e718bf	0	f	f	f	morsinska-negaz.webp
1de399d1-4951-4dfd-9495-6968c2f2240d	07060c7c	Найсмачніший бургер	naismacnisii-burger	Котлета з телятини, ціла голова камамберу на грилі, беконовий мармелад, трюфельний майо, помідор, ромен	370.00	0.375	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	naismacnisii-burger.webp
489f2581-a4fd-4a21-9aec-d9212f0faff1	f20fdbb9	Сет бургерів	set-burgeriv	MEAT, BEEF, CHICKEN,  SMASH BRO з однією котлетою,  картопля фрі і діпи, три соуси	999.00	1.400	524da3bb-bb79-4c53-ac60-a21dac89e202	0	f	f	f	set-burgeriv.webp
e1ca4e0c-ac8e-45cc-b238-f36e23bba849	b28e8780	Сет закусок	set-zakusok	Курячі крильця, ковбасне печиво, сирні палички, грінки з пармезаном, картопляні діпи, курячі нагетси, три соуси	720.00	1.100	524da3bb-bb79-4c53-ac60-a21dac89e202	0	f	f	f	set-zakusok.webp
397135c8-fc6c-408c-936b-4ffdeb0691f1	6dde4751	Сік на вибір	sik-na-vibir	Сік на вибір	50.00	0.300	525976cc-90e9-497d-9be3-2ce001e718bf	0	f	f	f	sik-na-vibir.webp
1b6f20ad-ad77-4350-bd42-a3ec487f37b7	680a86a4	Сирна ванна	sirna-vanna	Сир чедер, пармезан	170.00	0.160	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	sirna-vanna.webp
922814ee-01f8-4b38-a688-a2ae2bb3627b	0dc3d7b5	Сирний суп	sirnii-sup	Картопля, морква, цибуля, три види сиру арахіс, кріспі цибуля, крутони	140.00	0.270	b400c7ab-817b-47c0-a72b-a66597d440e2	0	f	f	f	sirnii-sup.webp
f00ef44f-425b-4a9f-bf8b-0089f52112b9	405f84dc	Сирно-часниковий	sirno-casnikovii	Сирно-чаниковий	30.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	sirno-casnikovii.webp
52912a48-e103-42ba-a64b-068a397b06ef	1fbacc7e	Smash Bro burger	smash-bro-burger	Подвійна смеш котлета з телятини, сир чеддер, смажений бекон, смажена цибуля, солоний огірок, помідор, халапеньйо, сирний соус з кімчі	370.00	0.400	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	smash-bro-burger.webp
753b49f7-f4c8-41db-812d-ea0753893b09	11e8bcf4	Сніданок холостяка	snidanok-xolostiaka	Картопля фрі, багато сирного соусу, бекон, ковбаска, перепелине яйце, кріспі цибуля, солоний огірок	330.00	0.330	aa5ace0d-00cb-4f37-a358-821d3afef2a1	0	f	f	f	snidanok-xolostiaka.webp
cc9bb588-04de-4697-b436-6414960334db	e46455d8	Spicy smash burger	spicy-smash-burger	Смеш котлета з телятини, шрірача майо, солоний огірок, цибуля кріспі, чеддер, беконовий мармелад	370.00	0.340	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	spicy-smash-burger.webp
d12b9757-903e-450d-8bc0-3926ec5247ed	642708a4	Surf & turf burger	surf-turf-burger	Котлета з телятини, тигрові креветки  на грилі, копчений майонез, сир чеддер, помідор, ромен	370.00	0.300	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	surf-turf-burger.webp
dc3e705b-2028-450f-9189-4fb373567837	0a03c09b	Тартар	tartar	Тартар	30.00	0.040	b1e398ad-5cbf-4f7d-8b55-16e1c6304d9a	0	f	f	f	tartar.webp
5a8415fd-7426-4e42-9c88-433152393295	e1171b78	Трюфельне фрі	triufelne-fri	Картопляні діпи, два трюфельних соуси, гриби ерінги, трюфельна паста, пармезан	275.00	0.260	aa5ace0d-00cb-4f37-a358-821d3afef2a1	0	f	f	f	triufelne-fri.webp
e2164431-b9a1-495c-8dac-6a9722e8e51e	74488064	Твій сет	tvii-set	MEAT, CHICKEN burger, курячі нагетси, тілапія в пивному клярі, картопляні діпи, три соуси	840.00	1.100	524da3bb-bb79-4c53-ac60-a21dac89e202	0	f	f	f	tvii-set.webp
b4117b7d-6ea7-417e-985a-94e24c6d4def	d70a2229	Батат фрі	batat-fri	Пармезан, соус трюфельний майо	220.00	0.170	aa5ace0d-00cb-4f37-a358-821d3afef2a1	0	f	f	f	batat-fri.webp
28cfa82a-4709-4787-8d42-a34972a949e0	43457f06	Beef burger	beef-burger	Котлета з телятини, яйце куряче, сир чеддер, помідор, солоний огірок, кріспі цибуля, трюфельний майо, ромен	340.00	0.400	f8a560e9-e0a5-4af3-bab2-5e2a1ac8a3de	0	f	f	f	beef-burger.webp
3363b324-7ecb-4797-aaa9-2641fa6cb75c	461ab0eb	Цезар з курячим філе	cezar-z-kuriacim-file	Мікс салатів, перепелині яйця, помідори черрі, куряче філе, крутони, пармезан, соус Цезар	240.00	0.260	b5241ce2-58f7-42e3-a053-19105384ddb6	0	f	f	f	cezar-z-kuriacim-file.webp
ab11263c-8c86-4a94-92dd-f3978939007d	9c861dc4	Crispy Wings	crispy-wings	Крила у паніровці, сирний соус, шрірача	220.00	0.300	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	crispy-wings.webp
32fb91cd-3d11-4ef8-af64-5a17b702ee8d	5ea60064	Джерки з курятини	dzerki-z-kuriatini	Карі, чилі, лайм	120.00	0.050	e3dabbf0-0e50-437d-bdd0-36e5b420082b	0	f	f	f	dzerki-z-kuriatini.webp
9f642192-d011-4f70-909f-53687f8206db	3f6053fb	Джерки зі свинини	dzerki-zi-svinini	Копчена паприка, коріандр, чорний перець	120.00	0.050	e3dabbf0-0e50-437d-bdd0-36e5b420082b	0	f	f	f	dzerki-zi-svinini.webp
07ce1ca1-0130-40d1-aa6b-e90da9d8f553	0aab0af6	Fish & Chips	fish-chips	Філе тілапії в пивному клярі, картопляні діпи, соус тартар	275.00	0.250	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	fish-chips.webp
9dadfcef-dcc9-4423-93c9-46d47564265d	9d6adf84	Фісташковий шніцель	fistaskovii-snicel	Кряче філе у паніровці, соус з фісташки та пармезану,сир горгонзола, картопляні кульки, в’ялені томати, сир пармезан	320.00	0.300	2ed614c9-5d7a-4d45-ac4f-78589dd79848	0	f	f	f	fistaskovii-snicel.webp
d72c42d7-6554-4f69-a2da-09626cd52467	7413b02f	Камамбер на грилі з трюфельним медом	kamamber-na-grili-z-triufelnim-medom	Камамбер на грилі, беконовий мармелад, трюфельний\nмед, яблучний чатні, м’ята, наша булочка на грилі	270.00	0.290	2ed614c9-5d7a-4d45-ac4f-78589dd79848	0	f	f	f	kamamber-na-grili-z-triufelnim-medom.webp
9b8cf3fa-fff9-4faa-99f4-0a88c3b19f7c	6b5f6539	Кільця кальмара в пивному клярі	kilcia-kalmara-v-pivnomu-kliari	Пивний кляр, соус копчений майонез	230.00	0.210	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	kilcia-kalmara-v-pivnomu-kliari.webp
26541893-c391-4474-89d3-f9668219b780	f30d49fa	Ковбаски гриль з овочами	kovbaski-gril-z-ovocami	Картопля по-селянськи, печериці і солодкий перець на грилі, соус BBQ	280.00	0.380	2ed614c9-5d7a-4d45-ac4f-78589dd79848	0	f	f	f	kovbaski-gril-z-ovocami.webp
9636629c-3bad-48c7-a305-ed435d061ee0	afa7cb3a	Ковбасне печиво	kovbasne-pecivo	Соус сирно-часниковий	160.00	0.180	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	kovbasne-pecivo.webp
560768a9-9c36-4e50-89f3-594937631e5e	c681e67d	Курячі крильця	kuriaci-krilcia	Соус огірковий майо	180.00	0.250	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	kuriaci-krilcia.webp
5768de14-3f41-4b9c-808e-4125bf4d446d	6e50b480	Курячі нагетси	kuriaci-nagetsi	Cоус сирно-часниковий	170.00	0.210	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	kuriaci-nagetsi.webp
483ffadc-cf8b-4eae-85a3-81b83782d24e	034161f2	Mac & Cheese з беконом	mac-cheese-z-bekonom	Паста, багато сиру, бекон, огірок солоний, кріспі цибуля	325.00	0.300	2ed614c9-5d7a-4d45-ac4f-78589dd79848	0	f	f	f	mac-cheese-z-bekonom.webp
d50388fe-8677-48d5-8146-7e7c1f649dd5	a30d37e5	Не проста яєчня	ne-prosta-iajecnia	Яєчня з двох яєць, бекон, ковбаски та помідор на грилі, фета, салат ромен, хліб	240.00	0.370	2ed614c9-5d7a-4d45-ac4f-78589dd79848	0	f	f	f	ne-prosta-iajecnia.webp
adff74f3-5f1b-4ee9-a6a2-a75a0d4eec58	78dd2ec8	Салат з фетою	salat-z-fetoiu	Салат ромен, тигрові креветки та апельсин на грилі, \nпомідори черрі, пармезан, цитрусова заправка	240.00	0.240	b5241ce2-58f7-42e3-a053-19105384ddb6	0	f	f	f	salat-z-fetoiu.webp
36cc33a1-2138-4833-8e64-518fdaf8eaa8	5ababd4f	Салат з креветками	salat-z-krevetkami	Салат ромен, тигрові креветки та апельсин на грилі, помідори черрі, пармезан, цитрусова заправка	295.00	0.220	b5241ce2-58f7-42e3-a053-19105384ddb6	0	f	f	f	salat-z-krevetkami.webp
5050c43e-24a4-4fee-a094-3a5120098617	d5ef1d3e	Сирні палички	sirni-palicki	Медово-гірчичний соус	175.00	0.210	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	sirni-palicki.webp
6a76056d-a3d2-4f19-8b4b-ddc63769f5c5	c4a4e501	Смажений халумі	smazenii-xalumi	Гострий мед	235.00	0.230	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	smazenii-xalumi.webp
ea62df67-a733-49e0-861f-7e3e90128e13	b3096c63	Тигрові креветки в панко	tigrovi-krevetki-v-panko	5 шт, пікантний крем сир, солодкий чилі	265.00	0.175	05362a6d-48e9-443b-a8c6-d405e7111a7c	0	f	f	f	tigrovi-krevetki-v-panko.webp
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: sanin
--

COPY public.schema_migrations (version, dirty) FROM stdin;
3	f
\.


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

\unrestrict fdWLD7pqMRumDVfmH3D68B60S48KJUL4gAIBKCxVtX6sEYjsg0qysD4nZfobPa4

