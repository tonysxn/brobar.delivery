"use client"

import * as React from "react"
import { StarRating } from "@/components/ui/star-rating"
import { CustomSwitch } from "@/components/ui/custom-switch"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { toast } from "sonner"
import { CheckCircle2 } from "lucide-react"

const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL;

export function ReviewForm() {
    const [foodRating, setFoodRating] = React.useState(0)
    const [serviceRating, setServiceRating] = React.useState(0)
    const [comment, setComment] = React.useState("")
    const [showContacts, setShowContacts] = React.useState(false)
    const [contactInfo, setContactInfo] = React.useState({
        name: "",
        email: "",
        phone: ""
    })
    const [isSubmitting, setIsSubmitting] = React.useState(false)
    const [isSuccess, setIsSuccess] = React.useState(false)

    // Regex for Ukrainian phone: +380XXXXXXXXX or 0XXXXXXXXX
    const phoneRegex = /^(\+380\d{9}|0\d{9})$/
    // Standard email regex
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    const isPhoneValid = contactInfo.phone.trim() === "" || phoneRegex.test(contactInfo.phone.trim())
    const isEmailValid = contactInfo.email.trim() === "" || emailRegex.test(contactInfo.email.trim())

    const hasAtLeastOneValidContact =
        (contactInfo.phone.trim() !== "" && isPhoneValid) ||
        (contactInfo.email.trim() !== "" && isEmailValid) ||
        contactInfo.name.trim() !== ""

    const areContactsValid = isPhoneValid && isEmailValid && hasAtLeastOneValidContact

    const isFormValid = foodRating > 0 && serviceRating > 0 &&
        (!showContacts || areContactsValid)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!isFormValid || isSubmitting) return

        setIsSubmitting(true)

        try {
            const response = await fetch(`${GATEWAY_URL}/reviews`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    food_rating: foodRating,
                    service_rating: serviceRating,
                    comment: comment,
                    phone: showContacts && contactInfo.phone ? contactInfo.phone : null,
                    email: showContacts && contactInfo.email ? contactInfo.email : null,
                    name: showContacts && contactInfo.name ? contactInfo.name : null,
                }),
            })

            if (!response.ok) {
                throw new Error("Failed to submit review")
            }

            setIsSuccess(true)
            toast.success("Дякуємо за ваш відгук!")
        } catch (error) {
            console.error("Error submitting review:", error)
            toast.error("Не вдалося надіслати відгук. Спробуйте ще раз.")
        } finally {
            setIsSubmitting(false)
        }
    }

    const inputClasses = "w-full bg-[#2a2a2a]/80 border border-white/10 rounded-md px-4 py-3 text-white placeholder:text-gray-500 focus:outline-none focus:border-yellow-500 transition-colors"

    if (isSuccess) {
        return (
            <div className="w-full max-w-2xl mx-auto flex flex-col gap-6 p-8 bg-zinc-950/50 backdrop-blur rounded-xl border border-white/5 items-center text-center">
                <CheckCircle2 className="w-16 h-16 text-green-500" />
                <h2 className="text-2xl font-bold text-white">Дякуємо за ваш відгук!</h2>
                <p className="text-gray-400">Ми цінуємо вашу думку та використаємо її для покращення нашого сервісу.</p>
            </div>
        )
    }

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-2xl mx-auto flex flex-col gap-8 p-6 bg-zinc-950/50 backdrop-blur rounded-xl border border-white/5">
            <div className="flex flex-col md:flex-row gap-8 justify-between">
                <StarRating
                    label="Страви"
                    value={foodRating}
                    onChange={setFoodRating}
                />
                <StarRating
                    label="Сервіс"
                    value={serviceRating}
                    onChange={setServiceRating}
                />
            </div>

            <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-gray-400 uppercase tracking-wider">Коментар</label>
                <textarea
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                    placeholder="Ваші враження..."
                    className={cn(inputClasses, "min-h-[120px] resize-y")}
                />
            </div>

            <div className="pt-2">
                <CustomSwitch
                    label="Будь ласка, залиште свої контакти"
                    checked={showContacts}
                    onCheckedChange={setShowContacts}
                />
            </div>

            {showContacts && (
                <div className="flex flex-col gap-4 animate-in fade-in slide-in-from-top-2 duration-300">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="flex flex-col gap-2">
                            <label className="text-sm font-medium text-gray-400 uppercase tracking-wider">Телефон</label>
                            <input
                                type="tel"
                                value={contactInfo.phone}
                                onChange={(e) => setContactInfo({ ...contactInfo, phone: e.target.value })}
                                placeholder="+380"
                                className={inputClasses}
                            />
                        </div>
                        <div className="flex flex-col gap-2">
                            <label className="text-sm font-medium text-gray-400 uppercase tracking-wider">Ел. пошта</label>
                            <input
                                type="email"
                                value={contactInfo.email}
                                onChange={(e) => setContactInfo({ ...contactInfo, email: e.target.value })}
                                placeholder="example@mail.com"
                                className={inputClasses}
                            />
                        </div>
                    </div>
                    <div className="flex flex-col gap-2">
                        <label className="text-sm font-medium text-gray-400 uppercase tracking-wider">Ім&apos;я (необов&apos;язково)</label>
                        <input
                            type="text"
                            value={contactInfo.name}
                            onChange={(e) => setContactInfo({ ...contactInfo, name: e.target.value })}
                            placeholder="Введіть ваше ім'я"
                            className={inputClasses}
                        />
                    </div>
                </div>
            )}

            <div className="pt-4">
                <Button
                    type="submit"
                    disabled={!isFormValid || isSubmitting}
                    className="w-full h-14 text-lg font-bold bg-yellow-500 hover:bg-yellow-400 text-black rounded-full shadow-lg shadow-yellow-500/10 transition-all active:scale-[0.98] disabled:opacity-50"
                >
                    {isSubmitting ? "Надсилання..." : "Надіслати відгук"}
                </Button>
            </div>
        </form>
    )
}
