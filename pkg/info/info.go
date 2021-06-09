package info

import (
	"time"
)

type Info struct {
	Name            string    `pkginfo:"pkgname" pkgdesc:"NAME" json:"name"`
	Base            string    `pkginfo:"pkgbase" pkgdesc:"BASE" json:"base"`
	Version         string    `pkginfo:"pkgver" pkgdesc:"VERSION" json:"version"`
	Description     string    `pkginfo:"pkgdesc" pkgdesc:"DESC" json:"description"`
	URL             string    `pkginfo:"url" pkgdesc:"URL" json:"url"`
	BuildDate       time.Time `pkginfo:"builddate" pkgdesc:"BUILDDATE" json:"build_date"`
	Packager        string    `pkginfo:"packager" pkgdesc:"PACKAGER" json:"packager"`
	Size            uint64    `pkginfo:"size" pkgdesc:"ISIZE" json:"installed_size"`
	Architecture    string    `pkginfo:"arch" pkgdesc:"ARCH" json:"architecture"`
	License         []string  `pkginfo:"license" pkgdesc:"LICENSE" json:"license"`
	Replaces        []string  `pkginfo:"replaces" pkgdesc:"REPLACES" json:"replaces"`
	Groups          []string  `pkginfo:"group" pkgdesc:"GROUPS" json:"groups"`
	Conflicts       []string  `pkginfo:"conflict" pkgdesc:"CONFLICTS" json:"conflicts"`
	Provides        []string  `pkginfo:"provides" pkgdesc:"PROVIDES" json:"provides"`
	Backup          []string  `pkginfo:"backup" json:"backup"`
	Depends         []string  `pkginfo:"depend" pkgdesc:"DEPENDS" json:"depends"`
	OptionalDepends []string  `pkginfo:"optdepend" pkgdesc:"OPTDEPENDS" json:"optional_depends"`
	MakeDepends     []string  `pkginfo:"makedepend" pkgdesc:"MAKEDEPENDS" json:"make_depends"`
	CheckDepends    []string  `pkginfo:"checkdepend" pkgdesc:"CHECKDEPENDS" json:"check_depends"`
}
