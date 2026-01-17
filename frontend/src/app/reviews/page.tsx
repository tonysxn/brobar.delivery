import Header from "@/components/header";
import { ReviewForm } from "@/components/review-form";
import Footer from "@/components/footer";

export default function ReviewsPage() {
    return (
        <div className="flex flex-col min-h-screen bg-transparent">
            <Header />

            <main className="flex-1 container mx-auto px-4 py-12 md:py-20 flex flex-col items-center justify-center">
                <div className="w-full max-w-4xl text-center mb-12">
                    <h1 className="text-4xl md:text-5xl font-bold mb-4 tracking-tight">
                        Ваш відгук
                    </h1>
                    <p className="text-gray-400 text-lg">
                        Ми цінуємо вашу думку та прагнемо ставати кращими для вас
                    </p>
                </div>

                <ReviewForm />
            </main>

            <Footer />
        </div>
    );
}
