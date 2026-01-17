import {toast} from "sonner"

export default function showErrorToast(message: string) {
    toast.error(message)
}

export function showSuccessToast(message: string) {
    toast.success(message)
}
