# swsc

`swsc` implements the Sliding-Window Site Characteristics (SWSC) method as described in https://doi.org/10.1093/molbev/msy069

Initial write was based on [PFinderUCE-SWSC-EN]

[PFinderUCE-SWSC-EN]: https://github.com/Tagliacollo/PFinderUCE-SWSC-EN

## Input

`swsc` reads a single nexus file processing two blocks: 

1. `DATA`, containing the UCE markers (unique by ID)
2. `SETS`, containing the UCE locations (unique by ID)

Example (any `...` denotes truncated content, see [PFinderUCE-SWSC-EN] for full file):

```
#NEXUS

BEGIN DATA;
DIMENSIONS  NTAX=10 NCHAR=5786;
FORMAT DATATYPE=DNA GAP=- MISSING=?;
MATRIX

sp1    AGAAAC...TGCAAAG
...
;

END;

BEGIN SETS;

    [loci]
    CHARSET chr_2828 = 1-376;
    CHARSET chr_4312 = 377-627;
    ...

     CHARPARTITION loci = 1:chr_2828, 2:chr_4312;

END;
```

## Output

`swsc` writes a single .csv file containing the chosen characteristics for each site of the UCEs. It can also produce a `cfg` for use by PartitionFinder2 by using the appropriate flag.

## Usage

Both `input` and `output` must be set, otherwise each flag is optional. See `swsc -help` for details.