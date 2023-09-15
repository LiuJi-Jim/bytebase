#!/bin/sh

echo "copying antlr4 sources"

mkdir -p src/plugins/sql-lsp/server/3rd-party/antlr4
cp -r node_modules/antlr4/src/antlr4/* src/plugins/sql-lsp/server/3rd-party/antlr4/