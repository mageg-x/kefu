package version

import "embed"

var Version string    // wukongim version
var Commit string     // git commit id
var CommitDate string // git commit date
var TreeState string  // git tree state

var ImWebFs embed.FS   // monitor静态资源
var KefuWebFS embed.FS // demo静态资源
