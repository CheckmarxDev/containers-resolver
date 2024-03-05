package syftExtractor

import (
	"fmt"
	"github.com/anchore/syft/syft/pkg"
)

func packageTypeToPackageManager(packageType pkg.Type) PackageManagerType {
	switch packageType {
	case pkg.ApkPkg, pkg.DebPkg, pkg.RpmPkg:
		return Oval
	case pkg.GemPkg:
		return Ruby
	case pkg.NpmPkg:
		return Npm
	case pkg.PythonPkg:
		return Python
	case pkg.PhpComposerPkg:
		return Php
	case pkg.JavaPkg:
		return Maven
	case pkg.GoModulePkg:
		return Go
	case pkg.DotnetPkg:
		return Nuget
	case pkg.CocoapodsPkg:
		return Ios
	case pkg.ConanPkg:
		return Cpp
	case pkg.JenkinsPluginPkg:
		return JenkinsPlugin
	case pkg.AlpmPkg:
		return Alpm
	case pkg.PortagePkg:
		return Portage
	case pkg.HackagePkg:
		return Hackage
	case pkg.RustPkg:
		return Rust
	case pkg.KbPkg:
		return Kb
	case pkg.DartPubPkg:
		return DartPub
	case pkg.Rpkg:
		return R
	case pkg.UnknownPkg, pkg.BinaryPkg, pkg.ErlangOTPPkg, pkg.GithubActionPkg, pkg.GithubActionWorkflowPkg,
		pkg.GraalVMNativeImagePkg, pkg.HexPkg, pkg.LinuxKernelPkg, pkg.LinuxKernelModulePkg,
		pkg.NixPkg, pkg.SwiftPkg, pkg.WordpressPluginPkg:
		return Unsupported
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
