package directory_path

import "path/filepath"

func Filepath_func() {
	//Dir 返回路径中除去最后一个路径元素的部分，即该路径最后一个元素所在的目录。在使用 Split 去掉最后一个元素后，会简化路径并去掉末尾的斜杠。
	//如果路径是空字符串，会返回 "."；如果路径由 1 到多个斜杠后跟 0 到多个非斜杠字符组成，会返回 "/"；其他任何情况下都不会返回以斜杠结尾的路径。
	//Base 函数返回路径的最后一个元素。在提取元素前会去掉末尾的斜杠。如果路径是 ""，会返回 "."；如果路径是只有一个斜杆构成的，会返回 "/"。
	var fp1 = filepath.Dir("./filepath_sample.go")
	var fp12 = filepath.Base("./filepath_sample.go")
	var fp2 = filepath.Dir("/root/my_dir/readme.txt")   //Dir 显示文件路径
	var fp22 = filepath.Base("/root/my_dir/readme.txt") //显示文件名

	println("fp1=", fp1)
	println("fp12=", fp12)
	println("fp2=", fp2)
	println("fp22=", fp22)

	//显示文件的扩展名
	println("readme.md Ext=", filepath.Ext("/root/readme.md"))
	println("/root/dir Ext=", filepath.Ext("/root/dir"))

	//相对路径和绝对路径
	println("./root/../my/txt.md isAbs=", filepath.IsAbs("./root/../my/txt.md"))
	println("/root/my/txt.md isAbs=", filepath.IsAbs("/root/my/txt.md"))

	absPath, _ := filepath.Abs("./root/../my/txt.md")
	println("./root/./my/txt.md Abs=", absPath)

	absPath, _ = filepath.Abs("/root/hello/./what/a/../my/txt.md")
	println("/root/hello/./what/a/../my/txt.md Abs=", absPath)
}
