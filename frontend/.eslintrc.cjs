module.exports = {
    "env": {
        "browser": true,
        "es2021": true
    },
    "extends": [
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended"
    ],
    "overrides": [
        {
            "files": ["*.svelte"],
            "processor": "svelte3/svelte3"
        },
        {
            "files": ["*.cy.ts"],
            "rules": {
                "@typescript-eslint/no-empty-function": "off",
                "@typescript-eslint/no-explicit-any": "off"
            }
        }
    ],
    "parser": "@typescript-eslint/parser",
    "parserOptions": {
        "ecmaVersion": "latest",
        "sourceType": "module"
    },
    "plugins": [
        "svelte3",
        "@typescript-eslint"
    ],
    "ignorePatterns": [
        "**/wailsjs/**/*",
        "**/cypress/support/**/*"
    ],
    "rules": {
    },
    "settings": {
        "svelte3/typescript": true
    },
    "globals": {
        "NodeJS": true
    }
}
