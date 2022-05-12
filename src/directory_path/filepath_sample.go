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

	//Rel 相对路径
	basePath := "/root/test"
	relPath, _ := filepath.Rel(basePath, "/root/test/ww/wang") //(根目录， 新的目录)-> 返回相对目录
	println("Relative path=", relPath)

	//根据操作系统的不同，分割Path环境变量。linux 是：， windows是;
	os_vars := filepath.SplitList("java_home/bin:other/bin")
	for _, v := range os_vars {
		print(v, "\t")
	}
	println()

	//dir 分割
	//Split函数将路径从最后一个路径分隔符后面位置分隔为两个部分（dir和file）并返回。如果路径中没有路径分隔符，函数返回值dir会设为空字符串，
	//file会设为path。两个返回值满足path == dir+file。
	println(filepath.Split("/root/books/stupid_rabbit.md"))

	//Join函数可以将任意数量的路径元素放入一个单一路径里，会根据需要添加路径分隔符。结果是经过简化的，所有的空字符串元素会被忽略。
	println("path joined", filepath.Join("root", "a", "b"))

	//FromSlash函数将path中的斜杠（'/'）替换为路径分隔符并返回替换结果
	println("from slash=", filepath.FromSlash("c:\\root\\hello\\.\\my\\..")) //转换成当前系统类型的分隔符
	println("to slash=", filepath.FromSlash("/root/hello/./my/.."))          //路径替换为linux风格

	//VolumeName函数返回最前面的卷名。如Windows系统里提供参数"C:\foo\bar"会返回"C:"；
	//Unix/linux系统的"\\host\share\foo"会返回"\\host\share"；其他平台会返回""。
	println("windows volume=", filepath.VolumeName("C:\\program\test"))
	println("linux volume=", filepath.VolumeName("/root/program"))

	//Clean函数通过单纯的词法操作返回和path代表同一地址的最短路径。
	//1. 将连续的多个路径分隔符替换为单个路径分隔符
	//2. 剔除每一个.路径名元素（代表当前目录）
	//3. 剔除每一个路径内的..路径名元素（代表父目录）和它前面的非..路径名元素
	//4. 剔除开始一个根路径的..路径名元素，即将路径开始处的"/.."替换为"/"（假设路径分隔符是'/'）
	println("clean path=", filepath.Clean("/root/.././hello/and/who/../../test/."))

}
