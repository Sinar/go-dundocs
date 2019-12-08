# go-dundocs
Tools to process questions + debates in State Assemblies (Selangor, Penang, and others)

Get the binaries from the Release tab for your environment (FreeBSD, OSX, Linux)

Example:

Running plan creates the detected Question boundaries:
```bash
$ ./dundocs plan --session=sesi1902 --source=./raw/Lisan/SOALAN\ MULUT\ \(1-20\).pdf
$ cat data/SOALAN\ MULUT\ \(1-20\)/split.yml 
  stateassemblysession: sesi1902
  hansardtype: 1
  hansardquestions:
  - questionnum: "1"
    pagenumstart: 1
    pagenumend: 7
  - questionnum: "2"
    pagenumstart: 8
    pagenumend: 10
  - questionnum: "3"
    pagenumstart: 11
    pagenumend: 11
  - questionnum: "4"
    pagenumstart: 12
    pagenumend: 13
  - questionnum: "5"
    pagenumstart: 14
    pagenumend: 20
  - questionnum: "6"
    pagenumstart: 21
    pagenumend: 21
  - questionnum: "7"
    pagenumstart: 22
    pagenumend: 23
  - questionnum: "8"
    pagenumstart: 24
    pagenumend: 24
  - questionnum: "9"
    pagenumstart: 25
    pagenumend: 26
  - questionnum: "10"
    pagenumstart: 27
    pagenumend: 27
  - questionnum: "11"
    pagenumstart: 28
    pagenumend: 30
  - questionnum: "12"
    pagenumstart: 31
    pagenumend: 34
  - questionnum: "13"
    pagenumstart: 35
    pagenumend: 35
  - questionnum: "14"
    pagenumstart: 36
    pagenumend: 38
  - questionnum: "15"
    pagenumstart: 39
    pagenumend: 39
  - questionnum: "16"
    pagenumstart: 40
    pagenumend: 43
  - questionnum: "17"
    pagenumstart: 44
    pagenumend: 45
  - questionnum: "18"
    pagenumstart: 46
    pagenumend: 46
  - questionnum: "19"
    pagenumstart: 47
    pagenumend: 48
  - questionnum: "20"
    pagenumstart: 49
    pagenumend: 50
```

Running split with the plan above ends with PDFs in splitout:
```bash
$ ./dundocs split --source=./raw/Lisan/SOALAN\ MULUT\ \(1-20\).pdf
$ ls splitout/SOALAN\ MULUT\ \(1-20\)/                
  sesi1902-soalan-1.pdf  sesi1902-soalan-16.pdf sesi1902-soalan-4.pdf
  sesi1902-soalan-10.pdf sesi1902-soalan-17.pdf sesi1902-soalan-5.pdf
  sesi1902-soalan-11.pdf sesi1902-soalan-18.pdf sesi1902-soalan-6.pdf
  sesi1902-soalan-12.pdf sesi1902-soalan-19.pdf sesi1902-soalan-7.pdf
  sesi1902-soalan-13.pdf sesi1902-soalan-2.pdf  sesi1902-soalan-8.pdf
  sesi1902-soalan-14.pdf sesi1902-soalan-20.pdf sesi1902-soalan-9.pdf
  sesi1902-soalan-15.pdf sesi1902-soalan-3.pdf
```

Full run for Lisan + BukanLisan in Selangor State Assembly:
```bash
$ ./dundocs plan --session=sesi1902 --source=./raw/Lisan/SOALAN\ MULUT\ \(21-40\).pdf
$ ./dundocs split --source=./raw/Lisan/SOALAN\ MULUT\ \(21-40\).pdf

$ ./dundocs plan --session=sesi1902 --source=./raw/Lisan/SOALAN\ MULUT\ \(41-60\).pdf
$ ./dundocs split --source=./raw/Lisan/SOALAN\ MULUT\ \(41-60\).pdf

$ ./dundocs plan --session=sesi1902 --source=./raw/Lisan/SOALAN\ MULUT\ \(61-80\).pdf
$ ./dundocs split --source=./raw/Lisan/SOALAN\ MULUT\ \(61-80\).pdf

$ ./dundocs plan --session=sesi1902 --source=./raw/Lisan/SOALAN\ MULUT\ \(81-100\).pdf
$ ./dundocs split --source=./raw/Lisan/SOALAN\ MULUT\ \(81-100\).pdf

$ ./dundocs plan --session=sesi1902 --source=./raw/BukanLisan/BukanLisan-1-20.pdf
$ ./dundocs split --source=./raw/BukanLisan/BukanLisan-1-20.pdf

$ ./dundocs plan --session=sesi1902 --source=./raw/BukanLisan/BukanLisan-21-40.pdf
$ ./dundocs split --session=sesi1902 --source=./raw/BukanLisan/BukanLisan-21-40.pdf

$ ./dundocs split --source=./raw/BukanLisan/BukanLisan-21-40.pdf

$ ./dundocs plan --session=sesi1902 --source=./raw/BukanLisan/BukanLisan-41-60.pdf
$ ./dundocs split --source=./raw/BukanLisan/BukanLisan-41-60.pdf
```
