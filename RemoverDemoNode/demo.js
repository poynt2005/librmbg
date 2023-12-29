var remover = require("./BackgroundRemover");
var path = require("path");
var fs = require("fs");
var mime = require("mime-types");
var readline = require("readline");

var rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

rl.question(
  "[Remover][Info ] Please enter the path name of image which you want to remove the background with\n",
  (answer) => {
    var srcFileAbsPath = path.resolve(answer.trim());

    if (!fs.existsSync(srcFileAbsPath)) {
      throw new Error(`file path you entered: ${srcFileAbsPath} is not exists`);
    }

    if (!mime.lookup(srcFileAbsPath).startsWith("image")) {
      throw new Error(
        `file path you entered: ${srcFileAbsPath} is not a valid image file`
      );
    }

    var srcFilePathParsed = path.parse(srcFileAbsPath);
    var destFilePath = path.join(
      path.dirname(srcFileAbsPath),
      srcFilePathParsed.name + "_no_bg" + srcFilePathParsed.ext
    );

    console.log(
      `[Remover][Info ] Destation image file will be written to ${destFilePath}`
    );
    rl.close();

    var bgRevmoer = new remover();
    var result = bgRevmoer.RemoveBackground(fs.readFileSync(srcFileAbsPath));
    fs.writeFileSync(destFilePath, result);
    bgRevmoer.Dispose();
  }
);
