# fakeData
Fabrication de dataset anonymisé pour Signaux Faibles

## Objectif: 
- Fabrication d'un dataset de test pour fournir une UI de test sans risquer de diffuser des informations sensibles

## Méthode d'anonymisation
- Remplacement des siret et des comptes urssaf avec des siret aléatoires
- supprimer des informations sirène à l'exception de la région et du code APE
- Application d'un coefficient aléatoire différent sur chaque type de données chiffrées (urssaf et financières)

## Utilisation
- Créer config.toml à partir de config.toml.example
- go run *.go
- profit