using BackgroundRemover;
using BackgroundRemover.Lib;

Console.WriteLine("[Remover][Info ] Please enter the path name of image which you want to remove the background with");

var srcImgPath = String.Empty;

if (args.Length > 0)
{
    srcImgPath = args[0];
}
else
{
    srcImgPath = Console.ReadLine() ?? throw new Exception("path name cannot be empty");
}
srcImgPath = Path.GetFullPath(srcImgPath);

if (!File.ReadAllBytes(srcImgPath).IsImage())
{
    throw new Exception($"target image path: {srcImgPath} is not a valid image");
}

var destImgPath = Path.Combine(Directory.GetParent(srcImgPath).FullName, Path.GetFileNameWithoutExtension(srcImgPath) + "_no_bg" + Path.GetExtension(srcImgPath));

Console.WriteLine($"[Remover][Info ] Destation image will write to {destImgPath}");


using BGRemover remover = new();
remover.RemoveBackground(srcImgPath, destImgPath);