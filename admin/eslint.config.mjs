import {dirname} from "path";
import {fileURLToPath} from "url";
import {FlatCompat} from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
    baseDirectory: __dirname,
});

export default [
    ...compat.extends(
        "next/core-web-vitals",
        "next",
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended"
    ),
    {
        rules: {
            "@typescript-eslint/ban-ts-comment": "off",
            "@typescript-eslint/explicit-module-boundary-types": "off",
            "@typescript-eslint/no-explicit-any": "off",
            "@typescript-eslint/ban-types": "off",
            "@typescript-eslint/no-unused-vars": "warn",
            "react/react-in-jsx-scope": "off",
        },
    },
];
