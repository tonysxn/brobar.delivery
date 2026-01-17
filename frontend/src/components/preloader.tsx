"use client";

import Image from "next/image";
import PreloaderImage from "@/resources/images/preloader.svg";
import { usePreloader } from "@/contexts/preloader-context";
const Preloader = () => {
    const { isLoading } = usePreloader();

    return (
        <div className={`preloader ${!isLoading ? "hidden" : ""}`}>
            <Image
                alt="preloader"
                src={PreloaderImage}
                width={200}
                className="pulse-loop"
            />
        </div>
    );
};

export default Preloader;
