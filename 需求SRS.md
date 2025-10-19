写一个golang程序，可以自动获取下载的opensdk下载地址，
需求如下：
1、需要分析 https://mirrors.tuna.tsinghua.edu.cn/Adoptium/ 网页，解析出<a href="25/" title="25">25/</a>的版本地址
2、进入每一个版本地址网页，再找到 <a href="jdk/" title="jdk">jdk/</a> 的opensdk的目录，
3、再进入jdk的网页，再找出所有的架构网页，例如：<a href="x64/" title="x64">x64/</a>等
4、再进入x64的网页，再找出所有的系统网页，例如：<a href="windows/" title="windows">windows/</a>
5、最后进入真正jdk下载地址网页，分析出所有的jdk版本下载地址。
6、把所有的jdk下载地址，保存成 jdkindex.json文件
[
  {
    "version": "openjdk-25",
    "url": "https://mirrors.huaweicloud.com/openjdk/25/openjdk-25_windows-x64_bin.zip"
  },
  {
    "version": "openjdk-24.0.2",
    "url": "https://mirrors.huaweicloud.com/openjdk/24.0.2/openjdk-24.0.2_windows-x64_bin.zip"
  }
 ]
7、生成Readme.md文档
8、所有生成的代码函数都需要有中文注释，涉及到调用第三个组件的函数，也需要有中文注释。
9、所有的生成代码，执行命令都自动执行，不需要问我。