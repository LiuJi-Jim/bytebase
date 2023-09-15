import antlr4 from 'antlr4';

export default class PlSqlLexerBase extends antlr4.Lexer {
  IsNewlineAtPos(pos: number) {
    const input = this.inputStream;
    const la = input.LA(pos);
    return la == -1 || String.fromCharCode(la) == '\n';
  }
}
