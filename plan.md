Sistemul de contorizare are urmatoarele caracteristici:
- Are elemente de sistem energetic (transforamatoare, consumatori, un singur separator, linii cu lungimi diferite)
- Putem sa presupunem ca are punct de plecare energia: o sursa de energie de 100W
- Are un sens energia
- Separatorul cand e inchis se schimba sensul circulatiei de puteri

Parametrii de bază ai fiecărui element din rețea:
- Transformatoare: rapoartele de transformare, pierderile și tipul (ridicător/coborâtor de tensiune).
- Consumatori: puterea activă și reactiva consumată, factorul de putere al fiecărui consumator.
- Separatorul: starea (închis/deschis), care va influența direcția de curgere a puterii și eventual poate adăuga impedanță.
- Linii de transmisie: rezistența și reactanța specifică pe unitatea de lungime și lungimea fiecărei linii.

Structura algoritmului:
- Inițializează parametrii de bază: impedanțele, factorii de putere ai consumatorilor, tensiunile nominale, etc.
- Stabilește direcția fluxului de putere: verifică starea separatorului.
- Calculează tensiunile și curenții prin metode nodale: folosește o soluție numerică iterativă pentru echilibrul de puteri.
- Contorizează și înregistrează: pentru fiecare nod, salvează puterea activă, curentul și tensiunea.
- Adaptează la schimbările de sarcină sau de configurație: actualizează rezultatele când intervine o schimbare în sistem.


Sa presupunem ca avem un sistem in felul urmator:
la intrate in sistem este sursa care debiteaza energie catre urmatoarele elemente:
- deoarece sursa de energie este de 100MW se va debita la o tensiune de 20kV si va intra energia intr-un transformator de 20/110
- o linie de 110kV de 10km care duce pana intr-o statie de transformare ce coboara tensiunea de la 110kV la 20kV pentru a alimenta primul consumator de 20MW si dupa alti 3km un alt consumator de 30MW, dupa o noua statie dar care urca tensiunea de la 20kV la 110kV
- o linie de 110kV de 20km care dupa pana intr-o statie de tansformare ce coboara tensiunea de la 110kV la 20kV pentru a alimenta un ultim consumator de 50MW

Explicație:
Structurile: Definim structuri pentru fiecare componentă (sursă, transformator, linie, consumator, separator) și pentru sistemul în ansamblu.
Calcul pierderi: Funcția calculateLineLosses estimează pierderile pe linii.
Iterații: Algoritmul traversează configurația sistemului, calculează pierderile și actualizează puterile disponibile.
Separator: Controlează dacă sursa adițională este inclusă în sistem.
Rulează acest cod pentru a modela sistemul inițial și extinde-l în funcție de cerințe.