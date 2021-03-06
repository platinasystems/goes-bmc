package main

import (
	"os"

	"github.com/platinasystems/goes/lang"
)

type patents string

const Patents patents = `
Additional IP Rights Grant (Patents)

"The Project" means the contents of this CLI source code repository.

"The Author" means the Project originator, Platina Systems, Inc.

"This implementation" means the copyrightable works distributed by the
Author as part of the Project.

The Author hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of the Project, where such license applies only to those
patent claims, both currently owned or controlled by the Author and
acquired in the future, licensable by the Author that are necessarily
infringed by this implementation of the Project.  This grant does not
include claims that would be infringed only as a consequence of further
modification of this implementation.  If you or your agent or exclusive
licensee institute or order or agree to the institution of patent
litigation against any entity (including a cross-claim or counterclaim
in a lawsuit) alleging that this implementation of the Project or any
code incorporated within this implementation of the Project constitutes
direct or contributory patent infringement, or inducement of patent
infringement, then any patent rights granted to you under this License
for this implementation of the Project shall terminate as of the date
such litigation is filed.
`

func (patents) String() string { return "patents" }
func (patents) Usage() string  { return "[show ]patents" }

func (patents) Apropos() lang.Alt {
	return lang.Alt{
		lang.EnUS: "print machine patents",
	}
}

func (s patents) Main(...string) error {
	os.Stdout.WriteString(string(s)[1:])
	return nil
}
