function getFilePath(file) {
  const fs = new FS("system");
  return fs.Abs(file);
}

// yao run scripts.test.csv
function csv() {
  return Process("plugins.docloader.text", getFilePath("test.csv"));
}
// yao run scripts.test.pdf
function pdf() {
  return Process("plugins.docloader.text", getFilePath("sample.pdf"));
}
// yao run scripts.test.pdf_password
function pdf_password() {
  return Process(
    "plugins.docloader.text",
    getFilePath("sample_password.pdf"),
    "password"
  );
}

// yao run scripts.test.html
function html() {
  return Process("plugins.docloader.text", getFilePath("test.html"));
}

// yao run scripts.test.txt
function txt() {
  return Process("plugins.docloader.text", getFilePath("test.txt"));
}

// yao run scripts.test.docx
function docx() {
  return Process("plugins.docloader.text", getFilePath("test.docx"));
}
// yao run scripts.test.xlsx
function xlsx() {
  return Process("plugins.docloader.text", getFilePath("test.xlsx"));
}
// yao run scripts.test.xlsx_password
function xlsx_password() {
  return Process(
    "plugins.docloader.text",
    getFilePath("test.xlsx"),
    "password"
  );
}
// yao run scripts.test.pptx
function pptx() {
  return Process("plugins.docloader.text", getFilePath("test.pptx"));
}

// yao run scripts.test.md
function md() {
  return Process("plugins.docloader.text", getFilePath("test.md"));
}

// yao run scripts.test.wiz
function wiz() {
  return Process("plugins.docloader.text", getFilePath("test1.ziw"));
}
