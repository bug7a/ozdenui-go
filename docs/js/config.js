/* Bismillah */

/*

OzdenUI for Go Website - Config - v26.11

The values of the site in one place: the brand, the links, the facts and the state of the module.
The site is the sister of the site of OzdenUI for Python (the same theme, header and build).

Started Date: November 2026
Webpage: https://bug7a.github.io/ozdenui-go/

Copyright 2026 Bugra Ozden <bugra.ozden@gmail.com> (https://bug7a.github.io/)
Licensed under the Apache License, Version 2.0
SPDX-License-Identifier: Apache-2.0

*/

"use strict";

const CONFIG = {

    // *** BRAND:
    familySelf: "go", // this site in the strip of the family (js/sections/family.js)
    brandName: "OzdenUI",
    logoFile: "assets/ozdenui-go-dark.svg", // the logo of OzdenUI FOR GO (342 x 132), its version for the dark bar

    // *** PATHS: (relative to this page)
    rootPath: "./",
    imagesPath: "images/",

    // *** THE MODULE:
    modulePath: "github.com/bug7a/ozdenui-go",
    version: "0.2.3",
    goVersion: "1.22",
    githubURL: "https://github.com/bug7a/ozdenui-go", // the repository of the module ("": the GitHub button is hidden)

    // *** LINKS:
    basicjsURL: "https://bug7a.github.io/basic.js/", // the site of basic.js
    // The handbook of the module: a page of this site (Ex: "handbook/"); its language follows the page.
    // "": there is none yet, and its buttons and links are hidden.
    handbookURL: "handbook/",
    componentsURL: "https://bug7a.github.io/js-components/",
    pythonURL: "https://bug7a.github.io/ozdenui-python/", // OzdenUI for Python, the twin of this module
    authorURL: "https://bug7a.github.io/",

    // *** FACTS: (of ozdenui-go v0.2.3: its README, kit_gen.go and web/)
    // WHY \u00A0: a number and its unit stay on one line
    componentCount: 83,           // the components of the kit, each a function with its own ...Props type (kit_gen.go)
    embedSize: "3.5\u00A0MB",     // basic.js and the components inside the program (go:embed of web/)
    license: "Apache 2.0",

    // *** THE SECTIONS OF THE HEADER:
    navKeys: ["code", "why", "components", "examples", "start"],

    // *** LANGUAGE:
    defaultLanguage: "en", // "tr", "en"
    languageStorageKey: "bjsgo_lang",

    // *** SEO:
    // The published address of this page (the docs/ folder of the repository, GitHub Pages).
    siteURL: "https://bug7a.github.io/ozdenui-go/",
    ogImage: "assets/og-image.jpg", // the link preview: the first screen of the page at 1200 x 630 (test/screenshot.js --og)

    // *** BEHAVIOR:
    rebuildOnResizeDelay: 180, // ms: the page is built again when the width changes

};
