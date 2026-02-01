--
-- PostgreSQL database dump
--

\restrict sBzt1qTnNVF2B3SKnzwcZEID4hOwOgdVgiyGizXFXR4v0YAP1M7HZDaRnhhfUCD

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
832d9904-a3b7-467c-9c19-7f00e4d7e157	Сети	sets	UtensilsCrossed	10
dd7129a7-c4e3-4164-8388-6adfe0a86942	Бургери	burgers	Sandwich	20
85d13f4c-0d5b-458f-a184-1259c221f684	Гарячі закуски	snacks	Drumstick	30
5d61c0dd-525f-4f3f-962b-ac048217d560	Гарячі страви	hot	ChefHat	40
3fe2c508-7b93-441e-b987-a4b62613f13a	Перші страви	persi%20stravi	Soup	50
2e3b506f-3f60-49e2-9d16-8280542b7341	Фрі та сир	fries-and-cheese	Utensils	60
5569857c-d66a-4c99-82b6-4b47e4f965f4	Салати	salad	Salad	70
96138456-761b-4894-8fd0-aa0b70ee7003	Холодні закуски	holodni	Beef	80
ba192bd7-b389-4044-bf62-cc36472ba688	Соуси	sauces	Milk	90
85bec573-7b1a-43eb-bb53-3eab24eb6303	Напої	drinks	Coffee	100
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
3de63581-4006-49f5-ac23-d6a0851bc555	2	Сет бургерів	сет-бургерів-2	MEAT, BEEF, CHICKEN,  SMASH BRO з однією котлетою,  картопля фрі і діпи, три соуси	999.00	1.400	832d9904-a3b7-467c-9c19-7f00e4d7e157	10	f	f	f	31fbe3ab-e728-4204-ab17-2e63169818d1.jpg
21cc41d6-637f-4fd8-b3a6-f3b9b702f269	1	Твій сет	твій-сет-1	MEAT, CHICKEN burger, курячі нагетси, тілапія в пивному клярі, картопляні діпи, три соуси	840.00	1.100	832d9904-a3b7-467c-9c19-7f00e4d7e157	20	f	f	f	e7d816bf-220a-41d6-b5d8-a073c338e766.jpg
86bd183d-238d-478c-9725-1710f7c8ef2b	3	Сет закусок	сет-закусок-3	Курячі крильця, ковбасне печиво, сирні палички, грінки з пармезаном, картопляні діпи, курячі нагетси, три соуси	720.00	1.100	832d9904-a3b7-467c-9c19-7f00e4d7e157	30	f	f	f	7cd1ee02-6172-4d71-b403-0a39c6093640.jpg
2feea981-0a56-4cf4-9457-c12df30fe96a	5	BIG BRO burger	big-bro-burger-5	Подвійна котлета з телятини, подвійний cир чеддер, цибулеві кільця, солоний огірок, помідор, копчений BBQ, ромен	470.00	0.600	dd7129a7-c4e3-4164-8388-6adfe0a86942	10	f	f	f	a85804b8-11e6-48d5-80b0-5860dfe03998.jpg
25e4b90a-750e-4db6-bd91-133e92bd113c	6	Meat burger	meat-burger-6	Котлета з телятини, свіжа цибуля, солоний огірок, помідор, сир чеддер, копчений майонез, ромен	340.00	0.360	dd7129a7-c4e3-4164-8388-6adfe0a86942	20	f	f	f	f56faff0-b1d0-42e8-bae5-d6a70c62d758.jpg
c2f991fe-5a3d-42ce-9573-1e9b9d2cab40	7	Chicken burger	chicken-burger-7	Куряча котлета, бекон, помідор,  йогуртовий соус, сир чеддер, ромен	340.00	0.380	dd7129a7-c4e3-4164-8388-6adfe0a86942	30	f	f	f	de5ece36-d9eb-499b-a36b-d458ba176eb6.jpg
778b0d31-e1fe-49f7-ba06-e281a135bb2a	8	Beef burger	beef-burger-8	Котлета з телятини, яйце куряче, сир чеддер, помідор, солоний огірок, кріспі цибуля, трюфельний майо, ромен	340.00	0.400	dd7129a7-c4e3-4164-8388-6adfe0a86942	40	f	f	f	ead34252-1b8f-4248-b867-c6ea043cdb39.jpg
9d57f836-ca0c-4bae-8cf8-25db10a8bd29	9	King Mushrooms burger	king-mushrooms-burger-9	Смеш котлета з телятини, два трюфельних  соуси, моцарела, солоний огірок, гриби ерінги	370.00	0.340	dd7129a7-c4e3-4164-8388-6adfe0a86942	50	f	f	f	824cde21-62e8-4fce-ad00-50eb9d7f6d64.jpg
b0391c7e-4782-4c4a-99c0-caa5a960aea3	10	Найсмачніший бургер	найсмачніший-бургер-10	Котлета з телятини, ціла голова камамберу на грилі, беконовий мармелад, трюфельний майо, помідор, ромен	370.00	0.375	dd7129a7-c4e3-4164-8388-6adfe0a86942	60	f	f	f	d4ca3e8c-f8fe-42a2-b43d-677be928fcd0.jpg
651335a4-8bbc-46db-bcc7-508dcca0e9e0	11	Spicy smash burger	spicy-smash-burger-11	Смеш котлета з телятини, шрірача майо, солоний огірок, цибуля кріспі, чеддер, беконовий мармелад	370.00	0.340	dd7129a7-c4e3-4164-8388-6adfe0a86942	70	f	f	f	8b4d0b8c-b284-4abf-8cab-7a50aaf14916.jpg
87f2eb0a-8323-4fcc-aadc-47573e2cba19	14	Surf & turf burger	surf--turf-burger-14	Котлета з телятини, тигрові креветки  на грилі, копчений майонез, сир чеддер, помідор, ромен	370.00	0.300	dd7129a7-c4e3-4164-8388-6adfe0a86942	80	f	f	f	4655e94a-9953-4a15-9786-be9ebe61fada.jpg
827b6062-2668-4425-b4e5-dd474113582d	15	Smash Bro burger	smash-bro-burger-15	Подвійна смеш котлета з телятини, сир чеддер, смажений бекон, смажена цибуля, солоний огірок, помідор, халапеньйо, сирний соус з кімчі	370.00	0.400	dd7129a7-c4e3-4164-8388-6adfe0a86942	90	f	f	f	9dd35d13-5411-4a0d-8072-ea2160ee9d32.jpg
9ff19986-2549-4e5a-b59f-db37f51f8734	17	Сирна ванна	сирна-ванна-17	Сир чедер, пармезан	170.00	0.160	dd7129a7-c4e3-4164-8388-6adfe0a86942	100	f	f	f	6c51ccab-636e-40b4-b45f-3702582d6cf3.jpg
18f0284a-9933-439e-827c-67aa967a55d4	18	Crispy Wings	crispy-wings-18	Крила у паніровці, сирний соус, шрірача	220.00	0.300	85d13f4c-0d5b-458f-a184-1259c221f684	10	f	f	f	a19e417b-5737-4abb-93ae-f2a71b5b98ec.jpg
74bed596-4069-41d2-8dca-1e0c2745c801	19	Fish & Chips	fish--chips-19	Філе тілапії в пивному клярі, картопляні діпи, соус тартар	275.00	0.250	85d13f4c-0d5b-458f-a184-1259c221f684	20	f	f	f	e3d8cd75-2295-448c-ac43-f4bb159a9ca7.jpg
58e98c28-a799-44f3-9ede-ca5f266a37d1	21	Кільця кальмара в пивному клярі	кільця-кальмара-в-пивному-клярі-21	Пивний кляр, соус копчений майонез	230.00	0.210	85d13f4c-0d5b-458f-a184-1259c221f684	30	f	f	f	a2ae83dc-1d23-4093-b451-e636678e0e97.jpg
e3baae85-a27c-420a-976e-ea7813633093	22	Ковбасне печиво	ковбасне-печиво-22	Соус сирно-часниковий	160.00	0.180	85d13f4c-0d5b-458f-a184-1259c221f684	40	f	f	f	9777ec83-0a27-4ccb-bda2-94231b2b0151.jpg
5f60056d-5c3e-43a3-ac1c-a0c8335c574a	23	Курячі крильця	курячі-крильця-23	Соус огірковий майо	180.00	0.250	85d13f4c-0d5b-458f-a184-1259c221f684	50	f	f	f	0f3db5b4-d56f-4b39-9a50-2191be81111c.jpg
fb2f9d12-252d-453f-b35e-c812ecb8f645	24	Курячі нагетси	курячі-нагетси-24	Cоус сирно-часниковий	170.00	0.210	85d13f4c-0d5b-458f-a184-1259c221f684	60	f	f	f	d1d63cef-62f9-4d0e-9973-41f22ffd760f.jpg
af3edbe4-2f70-4bb4-b1a6-f6051b6f0cd9	25	Сирні палички	сирні-палички-25	Медово-гірчичний соус	175.00	0.210	85d13f4c-0d5b-458f-a184-1259c221f684	70	f	f	f	596e0e46-569e-43a8-97c9-c3dd95d2c939.jpg
7d82012d-e312-4ff7-aea9-20f1eb54a44f	26	Смажений халумі	смажений-халумі-26	Гострий мед	235.00	0.230	85d13f4c-0d5b-458f-a184-1259c221f684	80	f	f	f	2fd42a1e-5fff-4b76-979c-36fd31bd0405.jpg
9f766b3f-c944-4bd9-8560-844302265b69	27	Тигрові креветки в панко	тигрові-креветки-в-панко-27	5 шт, пікантний крем сир, солодкий чилі	265.00	0.175	85d13f4c-0d5b-458f-a184-1259c221f684	90	f	f	f	0669826b-caf6-436f-a7c4-73ec52ec1de1.jpg
efdc9521-2380-4fe2-83dc-ca7ca405b133	28	Mac & Cheese з беконом	mac--cheese-з-беконом-28	Паста, багато сиру, бекон, огірок солоний, кріспі цибуля	325.00	0.300	5d61c0dd-525f-4f3f-962b-ac048217d560	10	f	f	f	3e255dd3-7e4c-41fc-ba11-32e16ac9b677.jpg
cc27d6bf-b4bf-4a23-84ef-88818c3d0b27	29	Камамбер на грилі з трюфельним медом	камамбер-на-грилі-з-трюфельним-медом-29	Камамбер на грилі, беконовий мармелад, трюфельний\r\nмед, яблучний чатні, м’ята, наша булочка на грилі	270.00	0.290	5d61c0dd-525f-4f3f-962b-ac048217d560	20	f	f	f	92156013-0e42-4e3a-a0fd-add5dcf85908.jpg
931e3839-c00d-4af2-a111-4587258ecea9	30	Ковбаски гриль з овочами	ковбаски-гриль-з-овочами-30	Картопля по-селянськи, печериці і солодкий перець на грилі, соус BBQ	280.00	0.380	5d61c0dd-525f-4f3f-962b-ac048217d560	30	f	f	f	1c29ece7-9bd7-49ce-9a4e-e0e7fc2aea4e.jpg
f0ac864e-4d00-48ba-a294-b2bd3c4b5c3c	31	Не проста яєчня	не-проста-яєчня-31	Яєчня з двох яєць, бекон, ковбаски та помідор на грилі, фета, салат ромен, хліб	240.00	0.370	5d61c0dd-525f-4f3f-962b-ac048217d560	40	f	f	f	f8b0de1e-6cf1-4dd5-86a7-cd66ed80e2f3.jpg
4a898f5d-0d97-41ae-b79a-266585053314	32	Фісташковий шніцель	фісташковий-шніцель-32	Кряче філе у паніровці, соус з фісташки та пармезану,сир горгонзола, картопляні кульки, в’ялені томати, сир пармезан	320.00	0.300	5d61c0dd-525f-4f3f-962b-ac048217d560	50	f	f	f	0832dbde-ad22-4ef6-9610-d51ffe4f34ef.jpg
97e880fd-9a8c-4a0a-b2ce-58655d68c451	34	Сирний суп	сирний-суп-34	Картопля, морква, цибуля, три види сиру арахіс, кріспі цибуля, крутони	140.00	0.270	3fe2c508-7b93-441e-b987-a4b62613f13a	10	f	f	f	25081709-d3c5-4bd8-98de-a3ed684a51f0.jpg
e11c5490-bedc-4af3-a716-b8e6a6b1c22a	38	Батат фрі	батат-фрі-38	Пармезан, соус трюфельний майо	220.00	0.170	2e3b506f-3f60-49e2-9d16-8280542b7341	10	f	f	f	da7d3852-3671-4bab-8e78-aefe580612bd.jpg
b079267e-fe9a-4d41-9be2-a1c741fc7a9f	39	Картопля фрі	картопля-фрі-39	Соус BBQ	120.00	0.170	2e3b506f-3f60-49e2-9d16-8280542b7341	20	f	f	f	9ea26f6a-4ea2-42c0-8290-63b6a9f8f9c6.jpg
b451c1a7-99b2-463e-b86d-7236e51a9356	40	Картопляні діпи	картопляні-діпи-40	Сирно-часниковий соус	140.00	0.170	2e3b506f-3f60-49e2-9d16-8280542b7341	30	f	f	f	96cbfbae-fc87-4b1f-820e-6b378935cc20.jpg
d2411719-0c15-461d-b99b-f41a35c2dc12	41	Сніданок холостяка	сніданок-холостяка-41	Картопля фрі, багато сирного соусу, бекон, ковбаска, перепелине яйце, кріспі цибуля, солоний огірок	330.00	0.330	2e3b506f-3f60-49e2-9d16-8280542b7341	40	f	f	f	cb104a25-d37c-40ec-83b2-7caa2752bc8b.jpg
325478f9-0f62-459c-a0c4-e02f5a3dcb59	42	Трюфельне фрі	трюфельне-фрі-42	Картопляні діпи, два трюфельних соуси, гриби ерінги, трюфельна паста, пармезан	275.00	0.260	2e3b506f-3f60-49e2-9d16-8280542b7341	50	f	f	f	f450b6a7-7b33-424a-8713-61c6b36f80db.jpg
f580f560-a2d4-4d81-82d0-c0631d5b83e4	43	Фрі з беконовим мармеладом	фрі-з-беконовим-мармеладом-43	Картопля фрі, сирний соус, беконовий мармелад, кріспі цибуля, трюфельний майо, пармезан	235.00	0.270	2e3b506f-3f60-49e2-9d16-8280542b7341	60	f	f	f	dc9a2e4c-c74b-4704-bf45-aefce6165268.jpg
bd793828-988e-430a-97dd-c3cc96bcd823	35	Салат з креветками	салат-з-креветками-35	Салат ромен, тигрові креветки та апельсин на грилі, помідори черрі, пармезан, цитрусова заправка	295.00	0.220	5569857c-d66a-4c99-82b6-4b47e4f965f4	10	f	f	f	8c7c02d3-5c22-4692-b41f-c9347f6bf30a.jpg
6bb34c25-e278-433b-b8c1-8c580d1500af	36	Салат з фетою	салат-з-фетою-36	Салат ромен, тигрові креветки та апельсин на грилі, \r\nпомідори черрі, пармезан, цитрусова заправка	240.00	0.240	5569857c-d66a-4c99-82b6-4b47e4f965f4	20	f	f	f	df1d199d-51b8-4e56-b213-44f91adf6ed5.jpg
a09079ab-6057-4dc2-bc6f-1876f4e89e86	37	Цезар з курячим філе	цезар-з-курячим-філе-37	Мікс салатів, перепелині яйця, помідори черрі, куряче філе, крутони, пармезан, соус Цезар	240.00	0.260	5569857c-d66a-4c99-82b6-4b47e4f965f4	30	f	f	f	3028ea6c-a719-4429-b901-9ca4ba141fbb.jpg
c8bf9989-1363-422c-9768-64cbf2927df8	44	Джерки з курятини	джерки-з-курятини-44	Карі, чилі, лайм	120.00	0.050	96138456-761b-4894-8fd0-aa0b70ee7003	10	f	f	f	6bd15057-9050-4b17-80d3-72858eeb8965.jpg
1c31eda5-8042-49f2-ae19-34615b61043c	45	Джерки зі свинини	джерки-зі-свинини-45	Копчена паприка, коріандр, чорний перець	120.00	0.050	96138456-761b-4894-8fd0-aa0b70ee7003	20	f	f	f	683afa6c-f975-4045-be9c-b7ad0e4d6e0f.jpg
87b5faa0-e78d-4cea-a041-6dcf2825134b	46	BBQ	bbq-46	BBQ	30.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	10	f	f	f	07156eb3-f7bd-4354-a7e4-fa4f383b00a3.jpg
a832e285-941d-407c-a15f-003d7228d46a	47	Гірчиця	гірчиця-47	Гірчиця	15.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	20	f	f	f	27e3d2bb-bf7b-406c-9237-13962a5c01f9.jpg
b957a79d-6392-4738-920b-a7705287b349	48	Кетчуп	кетчуп-48	Кетчуп	30.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	30	f	f	f	bc8ce1a7-7609-412b-bc70-76af84e6932a.jpg
f7c0aad0-8307-4c20-8653-73f4751837ef	49	Копчений BBQ	копчений-bbq-49	Копчений BBQ	40.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	40	f	f	f	13c2a138-5230-42b3-a669-da975e74d96f.jpg
3f6cecb8-8b52-4876-963c-2fb99856aab0	50	Копчений майонез	копчений-майонез-50	Копчений майонез	30.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	50	f	f	f	da8def1d-c75e-4c3d-83d8-1d29cd361032.jpg
06e178bd-96a7-4dfd-88c7-52d61b7ee2be	51	Майонез	майонез-51	Майонез	20.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	60	f	f	f	94ac1c04-57ee-4191-b2fb-b451ac80d9ee.jpg
5c573dfd-d501-4e61-adf7-5a9c9077351e	52	Медово-гірчичний	медово-гірчичний-52	Медово-гірчичний	30.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	70	f	f	f	f3a7c7c0-7a91-4846-a02a-517d9bcb2ab2.jpg
fd22d0f4-1de1-4200-b6cb-79f61f1f0593	53	Сирно-часниковий	сирно-часниковий-53	Сирно-чаниковий	30.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	80	f	f	f	68ef8044-859b-4bcc-bc71-168692328167.jpg
5229c4ef-73ee-4f4c-bdda-95c9f4e6e3c4	54	Тартар	тартар-54	Тартар	30.00	0.040	ba192bd7-b389-4044-bf62-cc36472ba688	90	f	f	f	ec4de4fd-3315-461a-b013-49dcec7ffa13.jpg
e848810a-17f9-46ee-8a8e-f15dc0a3763e	55	Clausthaler Б/А	clausthaler-ба-55	Clausthaler Б/А, Німеччина	100.00	0.500	85bec573-7b1a-43eb-bb53-3eab24eb6303	10	f	f	f	bc47fc71-fac9-4c13-a2aa-c7d7db39877b.jpg
716014a4-32e2-4ae8-8eb7-a7eae4703077	60	Сік на вибір	сік-на-вибір-60	Сік на вибір	50.00	0.300	85bec573-7b1a-43eb-bb53-3eab24eb6303	20	f	f	f	66c10017-6dc9-477e-a0d5-01b409e26777.jpg
22287532-cbcd-45a7-b454-847db10edcd8	57	Кока-кола	кока-кола-57	Кока-кола	60.00	0.500	85bec573-7b1a-43eb-bb53-3eab24eb6303	30	f	f	f	98285790-b616-4506-8cbf-7d445847af25.jpg
26878a4c-f248-43c5-8134-adbfa316c621	59	Моршинська негаз	моршинська-негаз-59	Негаз	40.00	0.500	85bec573-7b1a-43eb-bb53-3eab24eb6303	40	f	f	f	b9a5687b-31c1-48f7-ab50-29617d7df74f.jpg
22ec1b5c-b65a-4b6f-a97c-104c3bd29a7c	58	Моршинська газ	моршинська-газ-58	Газ	40.00	0.500	85bec573-7b1a-43eb-bb53-3eab24eb6303	50	f	f	f	806265a5-dbd8-4e2a-9589-3129d5b7b34e.jpg
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

\unrestrict sBzt1qTnNVF2B3SKnzwcZEID4hOwOgdVgiyGizXFXR4v0YAP1M7HZDaRnhhfUCD

