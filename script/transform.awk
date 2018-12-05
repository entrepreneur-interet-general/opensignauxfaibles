BEGIN { # Semi-column separated csv as input and output
  FS = ";"
  OFS = ";"
} 
NR==1 { # Change field names

  printf "%s", "\"Annee\""

  for (i=1; i<=NF; i++) {

    if ($i !~ "201"){ # Field without year
      f[++nf] = i
      printf "%s%s", OFS, $i
    } else { # Field with year
    match($i, "20..", year)
    field_name = gensub(" "year[0],"","g",$i) # Remove year from column name

    if (!visited[field_name]){
      nf++
      visited[field_name]++
      printf "%s%s", OFS, field_name;
    }
    f[nf , year[0]] = i 
  }
}
printf "%s", ORS
}
NR>1{
  for (current_y=2012; current_y<=2017; current_y++){ # FIX ME: years hardcoded
    printf "%i", current_y
    for (i=1; i<=nf; i++) {
      if (f[i])
        printf "%s%s", OFS, $(f[i]);
      else {
        # Only print fields relative to current year 
        printf "%s%s", OFS, $(f[i, current_y]);
      }
    }
    printf "%s", ORS # Each year on a new line
  }
}
