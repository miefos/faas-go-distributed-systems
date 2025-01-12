import sys
import importlib

def main():
    print("Im in execute.py :))")
    # Nome della funzione da eseguire
    function_name = sys.argv[1]

    # Importa il modulo function.py
    try:
        import function
    except ModuleNotFoundError:
        print("Errore: il file function.py non è stato trovato")
        sys.exit(1)

    # Controlla se la funzione esiste
    if not hasattr(function, function_name):
        print(f"Errore: la funzione '{function_name}' non è definita in function.py")
        sys.exit(1)

    # Esegue la funzione
    func = getattr(function, function_name)
    try:
        func()
    except Exception as e:
        print(f"Errore durante l'esecuzione della funzione: {e}")

if __name__ == "__main__":
    main()
