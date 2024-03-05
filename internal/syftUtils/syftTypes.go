package syftUtils

import (
	"fmt"
	"github.com/anchore/syft/syft/pkg"
)

func packageTypeToPackageManager(packageType pkg.Type) string {
	switch packageType {
	case pkg.ApkPkg, pkg.DebPkg, pkg.RpmPkg:
		return string(Oval)
	case pkg.GemPkg:
		return string(Ruby)
	case pkg.NpmPkg:
		return string(Npm)
	case pkg.PythonPkg:
		return string(Python)
	case pkg.PhpComposerPkg:
		return string(Php)
	case pkg.JavaPkg:
		return string(Maven)
	case pkg.GoModulePkg:
		return string(Go)
	case pkg.DotnetPkg:
		return string(Nuget)
	case pkg.CocoapodsPkg:
		return string(Ios)
	case pkg.ConanPkg:
		return string(Cpp)
	case pkg.JenkinsPluginPkg:
		return string(JenkinsPlugin)
	case pkg.AlpmPkg:
		return string(Alpm)
	case pkg.PortagePkg:
		return string(Portage)
	case pkg.HackagePkg:
		return string(Hackage)
	case pkg.RustPkg:
		return string(Rust)
	case pkg.KbPkg:
		return string(Kb)
	case pkg.DartPubPkg:
		return string(DartPub)
	case pkg.Rpkg:
		return string(R)
	case pkg.UnknownPkg, pkg.BinaryPkg, pkg.ErlangOTPPkg, pkg.GithubActionPkg, pkg.GithubActionWorkflowPkg,
		pkg.GraalVMNativeImagePkg, pkg.HexPkg, pkg.LinuxKernelPkg, pkg.LinuxKernelModulePkg,
		pkg.NixPkg, pkg.SwiftPkg, pkg.WordpressPluginPkg:
		return string(Unsupported)
	default:
		panic(fmt.Sprintf("Failed to cast syft package type: %s into SupportedPackageManagerPrefixType", packageType))
	}
}

type PackageManagerType string

const (
	Npm           PackageManagerType = "Npm"
	Nuget         PackageManagerType = "Nuget"
	Maven         PackageManagerType = "Maven"
	Python        PackageManagerType = "Python"
	Php           PackageManagerType = "Php"
	Ios           PackageManagerType = "Ios"
	Go            PackageManagerType = "Go"
	Cpp           PackageManagerType = "Cpp"
	Ruby          PackageManagerType = "Ruby"
	JenkinsPlugin PackageManagerType = "JenkinsPlugin"
	Alpm          PackageManagerType = "Alpm"
	Portage       PackageManagerType = "Portage"
	Hackage       PackageManagerType = "Hackage"
	Rust          PackageManagerType = "Rust"
	Kb            PackageManagerType = "Kb"
	DartPub       PackageManagerType = "DartPub"
	R             PackageManagerType = "R"
	Oval          PackageManagerType = "Oval"
	Unsupported   PackageManagerType = "Unsupported"
)
