"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { reviewsApi, Review } from "@/lib/api";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useState } from "react";
import { toast } from "sonner";
import { Loader2, Trash2, Star, Phone, Mail, User } from "lucide-react";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { AuthData } from "@/types/auth";

function StarDisplay({ rating }: { rating: number }) {
    return (
        <div className="flex items-center gap-0.5">
            {[1, 2, 3, 4, 5].map((star) => (
                <Star
                    key={star}
                    className={`w-4 h-4 ${star <= rating ? "fill-yellow-500 text-yellow-500" : "text-gray-300"}`}
                />
            ))}
        </div>
    );
}

export default function ReviewsPage() {
    const queryClient = useQueryClient();
    const [user] = useLocalStorage<AuthData | null>("auth", null);
    const [deletingId, setDeletingId] = useState<string | null>(null);

    const { data: reviews, isLoading } = useQuery({
        queryKey: ["reviews"],
        queryFn: reviewsApi.getAll,
    });

    const deleteMutation = useMutation({
        mutationFn: (id: string) => reviewsApi.delete(id, user?.access.token),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["reviews"] });
            toast.success("Відгук видалено");
            setDeletingId(null);
        },
        onError: (e: any) => {
            toast.error("Не вдалося видалити відгук: " + (e.response?.data?.error || e.message));
        },
    });

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleString("uk-UA", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-64">
                <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    return (
        <div className="p-6">
            <div className="flex items-center justify-between mb-6">
                <h1 className="text-2xl font-bold">Відгуки</h1>
                <span className="text-muted-foreground">
                    Всього: {reviews?.length || 0}
                </span>
            </div>

            {reviews && reviews.length > 0 ? (
                <div className="rounded-md border">
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead className="w-[180px]">Дата</TableHead>
                                <TableHead className="w-[100px]">Страви</TableHead>
                                <TableHead className="w-[100px]">Сервіс</TableHead>
                                <TableHead>Коментар</TableHead>
                                <TableHead className="w-[200px]">Контакти</TableHead>
                                <TableHead className="w-[80px]">Дії</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {reviews.map((review) => (
                                <TableRow key={review.id}>
                                    <TableCell className="text-sm text-muted-foreground">
                                        {formatDate(review.created_at)}
                                    </TableCell>
                                    <TableCell>
                                        <StarDisplay rating={review.food_rating} />
                                    </TableCell>
                                    <TableCell>
                                        <StarDisplay rating={review.service_rating} />
                                    </TableCell>
                                    <TableCell className="max-w-[300px]">
                                        <p className="truncate" title={review.comment}>
                                            {review.comment || <span className="text-muted-foreground italic">Без коментаря</span>}
                                        </p>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex flex-col gap-1 text-sm">
                                            {review.name && (
                                                <div className="flex items-center gap-1 text-muted-foreground">
                                                    <User className="w-3 h-3" />
                                                    <span>{review.name}</span>
                                                </div>
                                            )}
                                            {review.phone && (
                                                <div className="flex items-center gap-1 text-muted-foreground">
                                                    <Phone className="w-3 h-3" />
                                                    <span>{review.phone}</span>
                                                </div>
                                            )}
                                            {review.email && (
                                                <div className="flex items-center gap-1 text-muted-foreground">
                                                    <Mail className="w-3 h-3" />
                                                    <span>{review.email}</span>
                                                </div>
                                            )}
                                            {!review.name && !review.phone && !review.email && (
                                                <span className="text-muted-foreground italic">Анонімно</span>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <AlertDialog>
                                            <AlertDialogTrigger asChild>
                                                <Button
                                                    variant="ghost"
                                                    size="icon"
                                                    className="text-destructive hover:text-destructive"
                                                >
                                                    <Trash2 className="w-4 h-4" />
                                                </Button>
                                            </AlertDialogTrigger>
                                            <AlertDialogContent>
                                                <AlertDialogHeader>
                                                    <AlertDialogTitle>Видалити відгук?</AlertDialogTitle>
                                                    <AlertDialogDescription>
                                                        Цю дію неможливо скасувати. Відгук буде видалено назавжди.
                                                    </AlertDialogDescription>
                                                </AlertDialogHeader>
                                                <AlertDialogFooter>
                                                    <AlertDialogCancel>Скасувати</AlertDialogCancel>
                                                    <AlertDialogAction
                                                        onClick={() => deleteMutation.mutate(review.id)}
                                                        className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                                    >
                                                        Видалити
                                                    </AlertDialogAction>
                                                </AlertDialogFooter>
                                            </AlertDialogContent>
                                        </AlertDialog>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>
            ) : (
                <div className="flex flex-col items-center justify-center h-64 text-muted-foreground">
                    <p className="text-lg">Відгуків поки немає</p>
                    <p className="text-sm">Вони з&apos;являться тут після того, як клієнти їх залишать</p>
                </div>
            )}
        </div>
    );
}
