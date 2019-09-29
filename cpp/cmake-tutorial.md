原文 [cmake-tutorial](https://cmake.org/cmake-tutorial/)

以下是分步教程，涵盖了 `CMake` 帮助解决的常见构建系统用例。这些话题，已经在[精通CMake](http://www.kitware.com/products/books/CMakeBook.html)中作为独立的主题介绍过了，但是在示例工程中看一下他们是如何共同工作的是非常有帮助的。这个教程可以在 `CMake` 源码树的 [Tests/Tutorial](https://gitlab.kitware.com/cmake/cmake/blob/master/Help/guide/tutorial/index.rst) 目录中找到。每个步骤都有自己的子目录，其中包含该步骤的完整副本。

另请参阅 [cmake-buildsystem(7)](https://cmake.org/cmake/help/latest/manual/cmake-buildsystem.7.html#introduction) 和 [cmake-language(7)](https://cmake.org/cmake/help/latest/manual/cmake-language.7.html#organization) 手册页的简介部分，以获取CMake概念和源代码树组织的概述。

# 基本起点(步骤1)

最基本的项目是从源代码文件构建的可执行文件。对于简单项目，只需两行CMakeLists.txt文件就足够了。这将是我们教程的起点，该 `CMakeLists.txt` 看起来像：

```cmake
cmake_minimum_required (VERSION 2.6)
project (Tutorial)
add_executable(Tutorial tutorial.cxx)
```

注意，这个例子在 `CMakeLists.txt` 文件中使用小写命令。大写，小写，混合命令都被 CMake 所支持。源代码 `tutorial.cxx` 将计算数字的平方根，第一个版本非常简单，如下：

```cpp
//A simple program that computes the square root of a number
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
int main (int argc, char *argv[])
{
  if (argc < 2)
    {
    fprintf(stdout,"Usage: %s number\n",argv[0]);
    return 1;
    }
  double inputValue = atof(argv[1]);
  double outputValue = sqrt(inputValue);
  fprintf(stdout,"The square root of %g is %g\n",
          inputValue, outputValue);
  return 0;
}
```

## 添加版本号和配置的头文件

我们要添加的第一个特性是给我们的可执行文件和工程提供一个版本号。你可以只在源码中做这个工作，但是在 `CMakeLists.txt` 文件中做这些可以提供更多的灵活性。为了添加版本号，我们修改 `CMakeLists.txt` 文件如下：

```cmake
cmake_minimum_required (VERSION 2.6)
project (Tutorial)
# The version number.
set (Tutorial_VERSION_MAJOR 1)
set (Tutorial_VERSION_MINOR 0)
 
# configure a header file to pass some of the CMake settings
# to the source code
configure_file (
  "${PROJECT_SOURCE_DIR}/TutorialConfig.h.in"
  "${PROJECT_BINARY_DIR}/TutorialConfig.h"
  )
 
# add the binary tree to the search path for include files
# so that we will find TutorialConfig.h
include_directories("${PROJECT_BINARY_DIR}")
 
# add the executable
add_executable(Tutorial tutorial.cxx)
```

由于已配置的文件将写入二进制树，因此我们必须将该目录添加到路径列表中以便可以搜索 include 文件。然后我们在源码树中创建一个 `TutorialConfig.h.in` 文件，内容如下：

```c++
// the configured options and settings for Tutorial
#define Tutorial_VERSION_MAJOR @Tutorial_VERSION_MAJOR@
#define Tutorial_VERSION_MINOR @Tutorial_VERSION_MINOR@
```

当 CMake 配置这个头文件时，@Tutorial_VERSION_MAJOR@ 和  @Tutorial_VERSION_MINOR@ 的值将被 CMakeLists.txt 文件中的值替换。接下来，我们修改 `tutorial.cxx` 引入这个配置头文件来使用版本号。产生的源代码如下：

```c++
// A simple program that computes the square root of a number
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "TutorialConfig.h"
 
int main (int argc, char *argv[])
{
  if (argc < 2)
    {
    fprintf(stdout,"%s Version %d.%d\n",
            argv[0],
            Tutorial_VERSION_MAJOR,
            Tutorial_VERSION_MINOR);
    fprintf(stdout,"Usage: %s number\n",argv[0]);
    return 1;
    }
  double inputValue = atof(argv[1]);
  double outputValue = sqrt(inputValue);
  fprintf(stdout,"The square root of %g is %g\n",
          inputValue, outputValue);
  return 0;
}
```

主要的修改是包含 `TurorialConfig.h` 头文件并将版本号作为使用消息的一部分打印出来。

# 添加库(Adding a Library)(步骤2)

现在我们将库添加到项目中。这个库将包含我们自己实现的计算一个数字的平方根。然后可执行文件可以使用这个库而不是编译器提供的标准平方根函数。对于这个教程，我们将该库放到名为 `MathFunctions` 的子目录中。`CMakeLists.txt` 中将会有下面的一行：

```cmake
add_library(MathFunctions mysqrt.cxx)
```

源文件 `mysqrt.cxx` 有一个 `mysqrt` 函数，它提供了和编译器 sqrt 函数类似的功能。为了利用新库，我们在顶层`CMakeLists.txt` 文件中添加了 `add_subdirectory` 调用，以便构建该库。我们还添加了另一个 include 目录，以确保函数原型可以找到 `MathFunctions/MathFunctions.h` 头文件。最后一个修改的地方是为可执行文件添加了一个新的库。顶级 `CMakeLists.txt` 文件的第后几行看起来像：

```cmake
include_directories ("${PROJECT_SOURCE_DIR}/MathFunctions")
add_subdirectory (MathFunctions)

# add the executable
add_executable (Tutorial MathFunctions)
```

现在我们来考虑将 `MathFunctions` 库设为可选。在这个教程，真的没有任何原因来做这些，但是在有比较大的库或依赖第三方代码的库，你可能需要这样做。第一步是在顶级的 `CMakeLists.txt` 文件中添加一个选项。

```cmake
# should we use our own math functions?
option (USE_MYMATH "Use tutorial provided math implementation" ON)
```

这将在 CMake GUI 中显示，默认值为 ON，用户可以根据需要进行更改。这个设置将被存储在 cache 中以确保用户在这个项目中运行 CMake 时不需要每次都保持设置。下一个修改是构建和连接 `MathFunctions` 库的条件。我们修改顶级 `CMakeLists.txt` 的结尾几行来实现这些：

```cmake
# add the MathFunctions library?
#
if (USE_MYMATH)
  include_directories ("${PROJECT_SOURCE_DIR}/MathFunctions")
  add_subdirectory (MathFunctions)
  set (EXTRA_LIBS ${EXTRA_LIBS} MathFunctions)
endif (USE_MYMATH)
 
# add the executable
add_executable (Tutorial tutorial.cxx)
target_link_libraries (Tutorial  ${EXTRA_LIBS})
```

这里使用 USE_MYMATH 设置来决定是否编译和使用 MathFunctions。注意使用变量(在该例子中是 EXTRA_LIBS)来搜集所有可选库，以供后面连接到可执行库中。这是用于保持大型项目和许多可选组件整洁的常用方法。对源代码的相应更改相当简单：

```cpp
// A simple program that computes the square root of a number
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "TutorialConfig.h"
#ifdef USE_MYMATH
#include "MathFunctions.h"
#endif
 
int main (int argc, char *argv[])
{
  if (argc < 2)
    {
    fprintf(stdout,"%s Version %d.%d\n", argv[0],
            Tutorial_VERSION_MAJOR,
            Tutorial_VERSION_MINOR);
    fprintf(stdout,"Usage: %s number\n",argv[0]);
    return 1;
    }
 
  double inputValue = atof(argv[1]);
 
#ifdef USE_MYMATH
  double outputValue = mysqrt(inputValue);
#else
  double outputValue = sqrt(inputValue);
#endif
 
  fprintf(stdout,"The square root of %g is %g\n",
          inputValue, outputValue);
  return 0;
}
```

在源代码中， 我们也使用了 `USE_MYMATH` 。可以通过配置文件 `TutorialConfig.h.in` 将其从CMake提供给源代码，通过在配置文件中添加以下行：

```cpp
# cmakedefine USE_MYMATH
```

# 安装和测试(步骤3)

在该步骤中，我们给项目添加安装规则和测试支持。安装规则相当简单。对于 `MathFunctions` 库，我们通过将以下两行添加到 `MathFunctions` 的 `CMakeLists.txt` 文件，来设置待安装的库和头文件。

```cmake
install (TARGETS MathFunctions DESTINATION bin)
install (FILES MathFunctions.h DESTINATION include)
```

对于应用程序，将下面两行将添加到顶层 `CMakeLists.txt` 文件中来安装可执行文件和配置头文件：

```cmake
# add the install targets
install (TARGETS Tutorial DESTINATION bin)
install (FILES "${PROJECT_BINARY_DIR}/TutorialConfig.h"        
         DESTINATION include)
```

这就是全部。此时，您就可以构建该教程了。然后输入 make install（或者从IDE中构建 INSTALL 目标），然后它将安装合适的头文件，库，可执行文件。CMake 变量 `CMAKE_INSTALL_PREFIX` 被用来确定文件的安装根目录。增加测试也是一个相当简单的过程。在顶级 `CMakeLists.txt` 文件的后面，我们可以添加一些基本测试来验证应用程序是否正常运行。

```cmake
include(CTest)

# does the application run
add_test (TutorialRuns Tutorial 25)
# does it sqrt of 25
add_test (TutorialComp25 Tutorial 25)
set_tests_properties (TutorialComp25 PROPERTIES PASS_REGULAR_EXPRESSION "25 is 5")
# does it handle negative numbers
add_test (TutorialNegative Tutorial -25)
set_tests_properties (TutorialNegative PROPERTIES PASS_REGULAR_EXPRESSION "-25 is 0")
# does it handle small numbers
add_test (TutorialSmall Tutorial 0.0001)
set_tests_properties (TutorialSmall PROPERTIES PASS_REGULAR_EXPRESSION "0.0001 is 0.01")
# does the usage message work?
add_test (TutorialUsage Tutorial)
set_tests_properties (TutorialUsage PROPERTIES PASS_REGULAR_EXPRESSION "Usage:.*number")
```

在构建完成之后，可以运行 “ctest" 命令行工具来运行测试。第一个测试简单地验证应用正在运行，没有段错误或其他崩溃，程序返回零值。这是 CTest 测试的基本形式。下面的几个测试都使用了 `PASS_REGULAR_EXPRESSION` 测试忏悔验证测试的输出是否包含指定的字符串。在这个例子中，验证计算的平方根是否是正确的，并在提供错误参数时打印用法消息。如果你想添加更多的测试来测试不同的输入值，你可以考虑创建一个像下面这样的宏：

```cmake
#define a macro to simplify adding tests, then use it
macro (do_test arg result)
  add_test (TutorialComp${arg} Tutorial ${arg})
  set_tests_properties (TutorialComp${arg}
    PROPERTIES PASS_REGULAR_EXPRESSION ${result})
endmacro (do_test)
 
# do a bunch of result based tests
do_test (25 "25 is 5")
do_test (-25 "-25 is 0")
```

对于每次 do_test 调用，都会根据传递的参数将另一个测试和它的名字，输入和结果添加到项目中。

# 添加系统内省(步骤4)

接下来我们考虑给项目添加一些代码，这取决于目标平台可能不具备的功能。对于这个例子，我们将添加一些代码来决定目标平台是否有 log 和 exp 函数。当然，几乎所有的平台都有这些函数，但是对于该教程，我们假设它们并不常见。如果平台有 log，那么我们在 mysqrt 函数中使用它计算平方根。我们首先在顶级 `CMakeLists.txt` 文件中使用 `CheckFunctionExists.cmake 宏来测试这些函数的可用性。

```cmake
# does this system provide the log and exp functions?
include (CheckFunctionExists)
check_function exists (log HAVE_LOG)
check_function exists (exp HAVE_EXP)
```

接下来，如果 CMake 在平台上找到它们， 我们修改 `TutorialConfig.h.in` 定义的这些值。

```cpp
// does the platform provide exp and log functions?
#cmakedefine HAVE_LOG
#camkedefine HAVE_EXP
```

在 CMake 对 TurorialConfig.h 进行 `configure_file` 之前完成 log 和 exp 的测试是非常重要的。`configure_file` 命令使用 CMake 的当前设置立即配置这个文件。最后，如果 log 和 exp 在系统中是可用的话，我们可以在 mysprt 函数中使用下列代码提供一个基于它们的备用实现。

```cpp
// if we have both log and exp then use them
#if defined (HAVE_LOG) && defined (HAVE_EXP)
	return exp(log(x)*0.5);
#else // otherwise use an interative approach
...
```

# 添加生成文件和生成器(步骤5)

在该部分，我们将向您展示如何将生成的源文件添加到应用程序的构建过程中。对于这个例子，我们将创建一个预计算平方根的表作为构建处理的一部分，然后将该表编译进我们的应用程序中。为了完成它，我们首先需要一个生成表的程序。在 `MathFunctions` 子目录中，一个名为 `MakeTable.cxx` 的源文件将做这些。

```cpp
// A simple program that builds a sqrt table 
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
 
int main (int argc, char *argv[])
{
  int i;
  double result;
 
  // make sure we have enough arguments
  if (argc < 2)
    {
    return 1;
    }
  
  // open the output file
  FILE *fout = fopen(argv[1],"w");
  if (!fout)
    {
    return 1;
    }
  
  // create a source file with a table of square roots
  fprintf(fout,"double sqrtTable[] = {\n");
  for (i = 0; i < 10; ++i)
    {
    result = sqrt(static_cast<double>(i));
    fprintf(fout,"%g,\n",result);
    }
 
  // close the table with a zero
  fprintf(fout,"0};\n");
  fclose(fout);
  return 0;
}
```

注意，该表是作为有效的 C++ 代码生成的，并将输出的文件名作为参数传递。下一步是给 MathFunctions 的 CMakeLists.txt 文件添加合适的命令来构建 MakeTable 可执行文件，然后作为构建程序的一部分运行它。完成他只需要很少的命令，展示如下：

```cmake
# first we add the executable that generates the table
add_executable(MakeTable MakeTable.cxx)
 
# add the command to generate the source code
add_custom_command (
  OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/Table.h
  COMMAND MakeTable ${CMAKE_CURRENT_BINARY_DIR}/Table.h
  DEPENDS MakeTable
  )
 
# add the binary tree directory to the search path for 
# include files
include_directories( ${CMAKE_CURRENT_BINARY_DIR} )
 
# add the main library
add_library(MathFunctions mysqrt.cxx ${CMAKE_CURRENT_BINARY_DIR}/Table.h  )
```

首先，添加 MakeTable 可执行文件，就像添加任何其他可执行文件一样。然后，我们添加一个自定义命令，该命令指定如何运行 `MakeTable` 来生成  `Table.h` 。接下来，我们让 CMake 知道 `mysqrt.cxx` 依赖于生成的 `Table.h` 文件。这可以通过将生成的 `Table.h` 文件添加到 `MathFunctions` 库的源列表地来完成。我们还必须将当前二进制目录添加到 include 目录中，以确保 `Table.h` 可以被 `mysqrt.cxx` 找到并引入(include)。当这个工程被构建时，它将首先构建 `MakeTable` 可执行文件。然后它将运行 `MakeTable` 生成 `Table.h` 。最后，它将编译引入 `Table.h` 的 `mysqrt.cxx` 生成 `MathFunction` 库。此时，顶级 `CMakeLists.txt` 文件将包含我们已经添加的所有特性，如下：

```cmake
cmake_minimum_required (VERSION 2.6)
project (Tutorial)
include(CTest)
 
# The version number.
set (Tutorial_VERSION_MAJOR 1)
set (Tutorial_VERSION_MINOR 0)
 
# does this system provide the log and exp functions?
include (${CMAKE_ROOT}/Modules/CheckFunctionExists.cmake)
 
check_function_exists (log HAVE_LOG)
check_function_exists (exp HAVE_EXP)
 
# should we use our own math functions
option(USE_MYMATH 
  "Use tutorial provided math implementation" ON)
 
# configure a header file to pass some of the CMake settings
# to the source code
configure_file (
  "${PROJECT_SOURCE_DIR}/TutorialConfig.h.in"
  "${PROJECT_BINARY_DIR}/TutorialConfig.h"
  )
 
# add the binary tree to the search path for include files
# so that we will find TutorialConfig.h
include_directories ("${PROJECT_BINARY_DIR}")
 
# add the MathFunctions library?
if (USE_MYMATH)
  include_directories ("${PROJECT_SOURCE_DIR}/MathFunctions")
  add_subdirectory (MathFunctions)
  set (EXTRA_LIBS ${EXTRA_LIBS} MathFunctions)
endif (USE_MYMATH)
 
# add the executable
add_executable (Tutorial tutorial.cxx)
target_link_libraries (Tutorial  ${EXTRA_LIBS})
 
# add the install targets
install (TARGETS Tutorial DESTINATION bin)
install (FILES "${PROJECT_BINARY_DIR}/TutorialConfig.h"        
         DESTINATION include)
 
# does the application run
add_test (TutorialRuns Tutorial 25)
 
# does the usage message work?
add_test (TutorialUsage Tutorial)
set_tests_properties (TutorialUsage
  PROPERTIES 
  PASS_REGULAR_EXPRESSION "Usage:.*number"
  )
 
 
#define a macro to simplify adding tests
macro (do_test arg result)
  add_test (TutorialComp${arg} Tutorial ${arg})
  set_tests_properties (TutorialComp${arg}
    PROPERTIES PASS_REGULAR_EXPRESSION ${result}
    )
endmacro (do_test)
 
# do a bunch of result based tests
do_test (4 "4 is 2")
do_test (9 "9 is 3")
do_test (5 "5 is 2.236")
do_test (7 "7 is 2.645")
do_test (25 "25 is 5")
do_test (-25 "-25 is 0")
do_test (0.0001 "0.0001 is 0.01")
```

 `TutorialConfig.h.in` 如下：

```cpp
// the configured options and settings for Tutorial
#define Tutorial_VERSION_MAJOR @Tutorial_VERSION_MAJOR@
#define Tutorial_VERSION_MINOR @Tutorial_VERSION_MINOR@
#cmakedefine USE_MYMATH
 
// does the platform provide exp and log functions?
#cmakedefine HAVE_LOG
#cmakedefine HAVE_EXP
```

MathFunctions 的 CMakeLists.txt 文件看起来像：

```cmake
# first we add the executable that generates the table
add_executable(MakeTable MakeTable.cxx)
# add the command to generate the source code
add_custom_command (
  OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/Table.h
  DEPENDS MakeTable
  COMMAND MakeTable ${CMAKE_CURRENT_BINARY_DIR}/Table.h
  )
# add the binary tree directory to the search path 
# for include files
include_directories( ${CMAKE_CURRENT_BINARY_DIR} )
 
# add the main library
add_library(MathFunctions mysqrt.cxx ${CMAKE_CURRENT_BINARY_DIR}/Table.h)
 
install (TARGETS MathFunctions DESTINATION bin)
install (FILES MathFunctions.h DESTINATION include)
```

# 构建安装程序(Installer)(步骤6)

接下来，假设我们希望将我们的工程分发给其他人，以便他们可以使用它。我们希望为多种平台提供二进制和源码发布。这和我们之前在 [安装和测试(步骤3)](安装和测试(步骤3))所做的安装有一点儿不同，在安装和测试中我们安装了根据源代码构建的二进制文件。在这个示例中， 我们将构建支持二进制安装和软件包管理功能的安装包软件，如 cygwin，debian，RPM 等等。为此，我们将使用 `CPack` 特定于平台的安装程序，正如在章节 [使用 CPack包装](https://gitlab.kitware.com/cmake/community/wikis/doc/cpack/Packaging-With-CPack) 中所述。特别地，我们需要在顶级 `CMakeLists.txt` 文件的底部添加一些行。

```cmake
# build a CPack driven installer package
include (InstallRequiredSystemLibraries)
set (CPACK_RESOURCE_FILE_LICENSE  
     "${CMAKE_CURRENT_SOURCE_DIR}/License.txt")
set (CPACK_PACKAGE_VERSION_MAJOR "${Tutorial_VERSION_MAJOR}")
set (CPACK_PACKAGE_VERSION_MINOR "${Tutorial_VERSION_MINOR}")
include (CPack)
```

这就是全部，我们以引入 `InstallRequiredSystemLibraries` 作为开始。这个模块将引入该项目在当前平台所需要的所有运行库。接下来我们设置一些 `CPack` 变量来保存这个项目的许可证和版本信息。版本信息利用了我们在本教程前面设置的变量。最后，我们引入 `CPack` 模块，它将使用这些变量和一些设置安装程序需要的其他系统属性。

下一步是以常规方式构建项目，然后在其上运行 `CPack`。构建二进制版本，运行：

```shell
cpack --config CPackConfig.cmake
```

创建源码版本：

```shell
cpack --config CPackSourceConfig.cmake
```

# 添加 Dashboard(步骤7)

添加支持将测试结果提交到仪表盘是非常容易的。在本教程之前步骤中，我们已经为我们的项目定义了一些测试。我们仅需要运行这些测试并把它们提交到仪表盘上。为了包含支持仪表盘，我们在顶级 `CMakeLists.txt` 文件中引入 `CTest` 模块。

```cmake
# enable dashboard scripting
include (CTest)
```

我们也可以创建 *CTestConfig.cmake* 文件，来为仪表盘指定这个工程的名字

```cmake
set (CTEST_PROJECT_NAME "Tutorial")
```

`CTest` 会在运行时读取该文件。可以在项目中运行 CMake 来创建简单的仪表盘，修改二进制树的目录，然后运行 `ctest -D` 测试。你仪表盘的结果将上传到公共[仪表盘](https://open.cdash.org/index.php?project=PublicDashboard)。

