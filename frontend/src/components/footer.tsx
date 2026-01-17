import Link from "next/link";
import { Separator } from "@/components/ui/separator";

const Footer = () => {
    return (
        <footer
            className="background-dark w-full pt-3 md:pt-6 pb-3 flex flex-col items-center gap-3 text-white text-sm">
            <div className="flex flex-col md:flex-row gap-2 md:gap-10 items-center">
                <Link href="/privacy">Політика конфіденційності</Link>
                <Link href="/contract">Договір публічної оферти</Link>
            </div>
            <Separator className="hidden md:flex max-w-2/3 mx-auto" />
            <span>© 2026 Bro-Bar</span>
        </footer>
    )
}

export default Footer;
