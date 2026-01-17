"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { settingsApi, Setting } from "@/lib/api";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import { useState, useEffect } from "react";
import { toast } from "sonner";
import { Loader2, Pencil, Check, X } from "lucide-react";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { AuthData } from "@/types/auth";

export default function SettingsPage() {
    const queryClient = useQueryClient();
    const [user] = useLocalStorage<AuthData | null>("auth", null);
    const [editingSetting, setEditingSetting] = useState<Setting | null>(null);
    const [editValue, setEditValue] = useState("");
    const [jsonError, setJsonError] = useState<string | null>(null);
    const [boolValue, setBoolValue] = useState(false);

    const { data: settings, isLoading } = useQuery({
        queryKey: ["settings"],
        queryFn: settingsApi.getAll,
    });

    const mutation = useMutation({
        mutationFn: (variables: { key: string; value: string; type: string }) =>
            settingsApi.update(variables.key, variables.value, variables.type, user?.access.token),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["settings"] });
            setEditingSetting(null);
            toast.success("Setting updated successfully");
        },
        onError: (e: any) => {
            toast.error("Failed to update setting: " + (e.response?.data?.error || e.message));
        },
    });

    const handleEdit = (setting: Setting) => {
        setEditingSetting(setting);
        setJsonError(null);
        if (setting.setting_type === "boolean") {
            setBoolValue(setting.value === "true");
        } else if (setting.setting_type === "json") {
            try {
                const formatted = JSON.stringify(JSON.parse(setting.value), null, 2);
                setEditValue(formatted);
            } catch {
                setEditValue(setting.value);
            }
        } else {
            setEditValue(setting.value);
        }
    };

    const handleSave = () => {
        if (!editingSetting) return;

        let valueToSave = editValue;
        if (editingSetting.setting_type === "boolean") {
            valueToSave = boolValue ? "true" : "false";
        } else if (editingSetting.setting_type === "json") {
            try {
                JSON.parse(editValue);
            } catch (e) {
                toast.error("Invalid JSON format");
                return;
            }
        }

        mutation.mutate({
            key: editingSetting.key,
            value: valueToSave,
            type: editingSetting.setting_type,
        });
    };

    useEffect(() => {
        if (editingSetting?.setting_type === "json" && editValue) {
            try {
                JSON.parse(editValue);
                setJsonError(null);
            } catch (e) {
                setJsonError("Invalid JSON");
            }
        }
    }, [editValue, editingSetting?.setting_type]);

    if (isLoading) {
        return (
            <div className="flex items-center justify-center p-8">
                <Loader2 className="h-8 w-8 animate-spin" />
            </div>
        );
    }

    return (
        <div className="p-6 space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-2xl font-bold">Settings</h1>
            </div>

            <div className="border rounded-md">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Key</TableHead>
                            <TableHead>Type</TableHead>
                            <TableHead>Value</TableHead>
                            <TableHead className="w-[100px]">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {settings?.map((setting) => (
                            <TableRow key={setting.key}>
                                <TableCell className="font-medium">{setting.key}</TableCell>
                                <TableCell>
                                    <span className="text-xs text-muted-foreground">
                                        {setting.setting_type}
                                    </span>
                                </TableCell>
                                <TableCell className="max-w-[400px] truncate">
                                    {setting.setting_type === "boolean" ? (
                                        <span className={setting.value === "true" ? "text-green-600" : "text-muted-foreground"}>
                                            {setting.value === "true" ? "Enabled" : "Disabled"}
                                        </span>
                                    ) : (
                                        <span className="truncate block text-muted-foreground">{setting.value}</span>
                                    )}
                                </TableCell>
                                <TableCell>
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        onClick={() => handleEdit(setting)}
                                    >
                                        <Pencil className="h-4 w-4" />
                                    </Button>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>

            <Dialog open={!!editingSetting} onOpenChange={(open) => !open && setEditingSetting(null)}>
                <DialogContent className="sm:max-w-[600px]">
                    <DialogHeader>
                        <DialogTitle>Edit: {editingSetting?.key}</DialogTitle>
                    </DialogHeader>
                    <div className="py-4 space-y-4">
                        {editingSetting?.setting_type === "boolean" ? (
                            <div className="flex items-center space-x-3">
                                <Checkbox
                                    id="bool-value"
                                    checked={boolValue}
                                    onCheckedChange={(checked) => setBoolValue(checked === true)}
                                />
                                <Label htmlFor="bool-value" className="cursor-pointer">
                                    {boolValue ? "Enabled" : "Disabled"}
                                </Label>
                            </div>
                        ) : editingSetting?.setting_type === "json" ? (
                            <div className="space-y-2">
                                <div className="flex items-center justify-between">
                                    <Label>JSON Value</Label>
                                    {jsonError ? (
                                        <span className="text-xs text-destructive flex items-center gap-1">
                                            <X className="h-3 w-3" /> {jsonError}
                                        </span>
                                    ) : editValue && (
                                        <span className="text-xs text-green-600 flex items-center gap-1">
                                            <Check className="h-3 w-3" /> Valid JSON
                                        </span>
                                    )}
                                </div>
                                <Textarea
                                    value={editValue}
                                    onChange={(e) => setEditValue(e.target.value)}
                                    className={`font-mono min-h-[300px] text-sm ${jsonError ? "border-destructive" : ""}`}
                                    placeholder="Enter valid JSON..."
                                />
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => {
                                        try {
                                            const formatted = JSON.stringify(JSON.parse(editValue), null, 2);
                                            setEditValue(formatted);
                                        } catch { }
                                    }}
                                >
                                    Format JSON
                                </Button>
                            </div>
                        ) : (
                            <Input
                                value={editValue}
                                onChange={(e) => setEditValue(e.target.value)}
                            />
                        )}
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setEditingSetting(null)}>
                            Cancel
                        </Button>
                        <Button onClick={handleSave} disabled={mutation.isPending || !!jsonError}>
                            {mutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            Save
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}
