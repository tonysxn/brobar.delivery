import Header from "@/components/header";
import Footer from "@/components/footer";

export default function PrivacyPage() {
    return (
        <div>
            <Header />
            <main className="container mx-auto px-4 py-8 min-h-screen">
                <h1 className="text-3xl font-bold mb-6">Політика конфіденційності</h1>
                <div className="prose prose-invert max-w-none text-gray-300">
                    <h2 className="text-xl font-bold text-white mb-4">Шановний Клієнте!</h2>
                    <p className="mb-4">
                        Користуючись Сайтом «Brobar delivery» <a href="https://brobar.delivery" className="text-brand hover:underline">brobar.delivery</a> та вказуючи персональні дані, шляхом внесення їх у відповідні форми, або надаючи персональні дані під час розмови із Оператором, Ви надаєте згоду на обробку Ваших персональних даних на умовах що наведені нижче.
                    </p>
                    <p className="mb-4">
                        У випадку виникнення питань, пов’язаних із будь-якою дією або сукупністю дій стосовно Ваших персональних даних, Ви можете зв’язатися з нами, зателефонувавши на контактний номер телефону Brobar delivery <a href="tel:+380635009597" className="text-brand hover:underline">+38-(063)-500-95-97</a> або відправивши листа електронною поштою на адресу: <a href="mailto:info@brobar.delivery" className="text-brand hover:underline">info@brobar.delivery</a>.
                    </p>

                    <h3 className="text-lg font-bold text-white mt-6 mb-2">1. Загальні положення</h3>
                    <p className="mb-4">
                        Ця Політика конфіденційності (обробки персональних даних) (далі за текстом – Політика) визначає порядок обробки персональних даних та заходи щодо забезпечення їх безпеки на Сайті «Brobar delivery» brobar.delivery (далі за текстом – Сайт). Реалізація цієї Політики проводиться Адміністрацією Сайту та уповноваженими нею особами (далі за текстом – Оператор).
                    </p>
                    <p className="mb-4">
                        Оператор у своїй діяльності керується нормами чинного національного законодавства України, міжнародними стандартами та вимогами, й визначає своєю провідною метою й головною умовою здійснення своєї діяльності – дотримання прав і свобод людини і громадянина під час обробки персональних даних, у тому числі, однак не обмежуючись, й захисту прав на недоторканність приватного життя, особисту і сімейну таємницю.
                    </p>
                    <p className="mb-4">
                        Ця Політика застосовується до всієї інформації, яку Оператор може отримати про відвідувачів Сайту «Brobar delivery» brobar.delivery та клієнтів доставки «Brobar delivery» відповідно по положень, викладених у цій Політиці та нормах чинного законодавства.
                    </p>
                </div>
            </main>
            <Footer />
        </div>
    );
}
