"use client";

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Order, ordersApi } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";

interface OrderDetailsDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    orderId: string | null;
    token?: string;
}

const STATUS_MAP: Record<string, { label: string; color: string }> = {
    "pending": { label: "New", color: "bg-blue-500" },
    "paid": { label: "Paid", color: "bg-green-500" },
    "shipping": { label: "Delivering", color: "bg-purple-500" },
    "completed": { label: "Completed", color: "bg-gray-500" },
    "cancelled": { label: "Cancelled", color: "bg-red-500" },
    // Kept for future reference
    "confirmed": { label: "Confirmed", color: "bg-green-600" },
    "cooking": { label: "Cooking", color: "bg-yellow-500" },
    "ready": { label: "Ready", color: "bg-orange-500" },
    "delivery": { label: "Delivering", color: "bg-purple-500" },
};

function StatusBadge({ statusId }: { statusId: string }) {
    const status = STATUS_MAP[statusId] || { label: "Unknown", color: "bg-gray-400" };
    return <Badge className={`${status.color} hover:${status.color}`}>{status.label}</Badge>;
}

export function OrderDetailsDialog({ open, onOpenChange, orderId, token }: OrderDetailsDialogProps) {
    const { data: order, isLoading } = useQuery({
        queryKey: ["order", orderId],
        queryFn: () => (orderId ? ordersApi.getOne(orderId, token) : Promise.reject("No ID")),
        enabled: !!orderId && open,
    });

    const formatDate = (dateString?: string) => {
        if (!dateString) return "-";
        return new Date(dateString).toLocaleString("uk-UA", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    const formatPrice = (price: number) => `${price} ₴`;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-3xl max-h-[90vh] flex flex-col">
                <DialogHeader>
                    <DialogTitle>Order Details #{orderId?.slice(0, 8)}</DialogTitle>
                </DialogHeader>

                {isLoading ? (
                    <div className="flex items-center justify-center p-8">
                        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                    </div>
                ) : order ? (
                    <ScrollArea className="flex-1 pr-4">
                        <div className="space-y-6">
                            {/* Header Info */}
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 p-4 border rounded-lg bg-card text-card-foreground">
                                <div>
                                    <p className="text-xs text-muted-foreground">Status</p>
                                    <StatusBadge statusId={order.status_id} />
                                </div>
                                <div>
                                    <p className="text-xs text-muted-foreground">Date</p>
                                    <p className="font-medium">{formatDate(order.created_at)}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-muted-foreground">Total</p>
                                    <p className="font-medium text-lg">{formatPrice(order.total_price)}</p>
                                </div>
                                <div>
                                    <p className="text-xs text-muted-foreground">Payment</p>
                                    <p className="font-medium capitalize">{order.payment_method}</p>
                                </div>
                            </div>

                            {/* Customer & Delivery */}
                            <div className="grid md:grid-cols-2 gap-6">
                                <div className="space-y-4">
                                    <h3 className="font-semibold text-lg border-b pb-2">Customer</h3>
                                    <div className="space-y-2 text-sm">
                                        <div className="grid grid-cols-[80px_1fr]">
                                            <span className="text-muted-foreground">Name:</span>
                                            <span className="font-medium">{order.name}</span>
                                        </div>
                                        <div className="grid grid-cols-[80px_1fr]">
                                            <span className="text-muted-foreground">Phone:</span>
                                            <span className="font-medium">{order.phone}</span>
                                        </div>
                                        {/* Add Email if available in API */}
                                    </div>
                                </div>

                                <div className="space-y-4">
                                    <h3 className="font-semibold text-lg border-b pb-2">Delivery</h3>
                                    <div className="space-y-2 text-sm">
                                        <div className="grid grid-cols-[80px_1fr]">
                                            <span className="text-muted-foreground">Type:</span>
                                            <span className="font-medium capitalize">{order.delivery_type_id}</span>
                                        </div>
                                        {order.delivery_type_id === "delivery" && (
                                            <>
                                                <div className="grid grid-cols-[80px_1fr]">
                                                    <span className="text-muted-foreground">Address:</span>
                                                    <span>{order.address}</span>
                                                </div>
                                                {order.zone && (
                                                    <div className="grid grid-cols-[80px_1fr]">
                                                        <span className="text-muted-foreground">Zone:</span>
                                                        <span>{order.zone}</span>
                                                    </div>
                                                )}
                                                {order.delivery_door && (
                                                    <div className="grid grid-cols-[80px_1fr] text-blue-600 font-medium">
                                                        <span className="text-muted-foreground">Service:</span>
                                                        <span>
                                                            Door Delivery
                                                            {order.delivery_door_price ? ` (+${formatPrice(order.delivery_door_price)})` : ""}
                                                        </span>
                                                    </div>
                                                )}
                                            </>
                                        )}
                                    </div>
                                </div>
                            </div>

                            {/* Wishes - Full Width */}
                            {order.wishes && (
                                <div className="space-y-2 border-t pt-4">
                                    <h3 className="font-semibold text-lg">Wishes</h3>
                                    <div className="p-3 bg-muted rounded-md text-sm italic">
                                        {order.wishes}
                                    </div>
                                </div>
                            )}

                            {/* Order Items */}
                            <div className="space-y-2">
                                <h3 className="font-semibold text-lg border-b pb-2">Items</h3>
                                <div className="border rounded-md">
                                    <Table>
                                        <TableHeader>
                                            <TableRow>
                                                <TableHead>Product</TableHead>
                                                <TableHead>Price</TableHead>
                                                <TableHead className="text-center">Qty</TableHead>
                                                <TableHead className="text-right">Total</TableHead>
                                            </TableRow>
                                        </TableHeader>
                                        <TableBody>
                                            {order.items?.map((item) => (
                                                <TableRow key={item.id}>
                                                    <TableCell>
                                                        <div className="font-medium">{item.name}</div>
                                                        {(item.product_variation_name || item.product_variation_group_name) && (
                                                            <div className="text-xs text-muted-foreground">
                                                                {item.product_variation_group_name}: {item.product_variation_name}
                                                            </div>
                                                        )}
                                                    </TableCell>
                                                    <TableCell>{formatPrice(item.price)}</TableCell>
                                                    <TableCell className="text-center">{item.quantity}</TableCell>
                                                    <TableCell className="text-right font-medium">
                                                        {formatPrice(item.total_price)}
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                            <TableRow>
                                                <TableCell colSpan={3} className="text-right font-semibold">
                                                    Delivery {order.delivery_door_price ? "(Base)" : ""}
                                                </TableCell>
                                                <TableCell className="text-right">
                                                    {formatPrice(order.delivery_cost - (order.delivery_door_price || 0))}
                                                </TableCell>
                                            </TableRow>
                                            {order.delivery_door_price > 0 && (
                                                <TableRow>
                                                    <TableCell colSpan={3} className="text-right font-semibold">
                                                        Delivery (Door)
                                                    </TableCell>
                                                    <TableCell className="text-right">
                                                        {formatPrice(order.delivery_door_price)}
                                                    </TableCell>
                                                </TableRow>
                                            )}
                                            <TableRow>
                                                <TableCell colSpan={3} className="text-right font-bold">
                                                    TOTAL
                                                </TableCell>
                                                <TableCell className="text-right font-bold text-lg">
                                                    {formatPrice(order.total_price)}
                                                </TableCell>
                                            </TableRow>
                                        </TableBody>
                                    </Table>
                                </div>
                            </div>
                        </div>
                    </ScrollArea>
                ) : (
                    <div className="p-8 text-center text-muted-foreground">Failed to load order details</div>
                )}
            </DialogContent>
        </Dialog>
    );
}
