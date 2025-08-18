#!/usr/bin/env python3
"""
Protein structure parser using BioPython
Parses PDB, CIF, and mmCIF files and extracts structural information
"""

import sys
import json
import warnings
from pathlib import Path
from collections import defaultdict

# Suppress BioPython warnings
warnings.filterwarnings('ignore')

try:
    from Bio import PDB
    from Bio.PDB.MMCIF2Dict import MMCIF2Dict
except ImportError:
    print(json.dumps({"error": "BioPython not installed"}))
    sys.exit(1)


class ProteinParser:
    """Parser for protein structure files"""
    
    def __init__(self, filepath):
        self.filepath = Path(filepath)
        self.structure = None
        self.parser = None
        
    def parse(self):
        """Parse the structure file based on extension"""
        ext = self.filepath.suffix.lower()
        
        if ext == '.pdb':
            self.parser = PDB.PDBParser(QUIET=True)
        elif ext in ['.cif', '.mmcif']:
            self.parser = PDB.MMCIFParser(QUIET=True)
        else:
            raise ValueError(f"Unsupported file format: {ext}")
            
        self.structure = self.parser.get_structure('protein', str(self.filepath))
        
    def extract_chains(self):
        """Extract chain information"""
        chains = []
        
        for model in self.structure:
            for chain in model:
                chain_type = self._determine_chain_type(chain)
                residue_count = len([r for r in chain if r.id[0] == ' '])
                atom_count = sum(len(r) for r in chain)
                
                chains.append({
                    "id": chain.id,
                    "residues": residue_count,
                    "atoms": atom_count,
                    "type": chain_type
                })
                
        return chains
    
    def _determine_chain_type(self, chain):
        """Determine the type of chain (Protein, DNA, RNA, Water)"""
        residue_types = set()
        
        for residue in chain:
            resname = residue.resname.strip()
            
            # Check for water
            if resname in ['HOH', 'WAT', 'H2O']:
                residue_types.add('Water')
            # Check for DNA
            elif resname in ['DA', 'DC', 'DG', 'DT', 'DU']:
                residue_types.add('DNA')
            # Check for RNA
            elif resname in ['A', 'C', 'G', 'U', 'I']:
                residue_types.add('RNA')
            # Check for standard amino acids
            elif resname in PDB.Polypeptide.standard_aa_names:
                residue_types.add('Protein')
                
        # Return the most common type
        if 'Protein' in residue_types:
            return 'Protein'
        elif 'DNA' in residue_types:
            return 'DNA'
        elif 'RNA' in residue_types:
            return 'RNA'
        elif 'Water' in residue_types:
            return 'Water'
        else:
            return 'Other'
    
    def extract_ligands(self):
        """Extract ligand information"""
        ligands = []
        ligand_dict = defaultdict(lambda: {"atoms": 0, "sites": []})
        
        for model in self.structure:
            for chain in model:
                for residue in chain:
                    # Check if it's a hetero residue (ligand)
                    if residue.id[0] != ' ':
                        resname = residue.resname.strip()
                        # Skip water molecules
                        if resname not in ['HOH', 'WAT', 'H2O']:
                            atom_count = len(residue)
                            site = f"Chain {chain.id}:{residue.id[1]}"
                            
                            ligand_dict[resname]["atoms"] += atom_count
                            ligand_dict[resname]["sites"].append(site)
        
        for name, info in ligand_dict.items():
            ligands.append({
                "name": name,
                "atoms": info["atoms"],
                "binding_site": ", ".join(set(info["sites"][:3]))  # Show first 3 sites
            })
            
        return ligands
    
    def extract_domains(self):
        """Extract domain information (simplified)"""
        # Note: Actual domain detection requires additional annotation data
        # This is a simplified version based on chain segments
        domains = []
        
        for model in self.structure:
            for chain in model:
                residues = [r for r in chain if r.id[0] == ' ']
                if len(residues) > 50:
                    # Create artificial domains for demonstration
                    mid_point = len(residues) // 2
                    
                    domains.append({
                        "name": f"N-terminal domain (Chain {chain.id})",
                        "range": f"1-{mid_point}",
                        "chains": [chain.id],
                        "atoms": sum(len(r) for r in residues[:mid_point])
                    })
                    
                    domains.append({
                        "name": f"C-terminal domain (Chain {chain.id})",
                        "range": f"{mid_point+1}-{len(residues)}",
                        "chains": [chain.id],
                        "atoms": sum(len(r) for r in residues[mid_point:])
                    })
                    
        return domains
    
    def get_total_atoms(self):
        """Get total number of atoms in the structure"""
        total = 0
        for model in self.structure:
            for chain in model:
                for residue in chain:
                    total += len(residue)
        return total
    
    def to_json(self):
        """Convert structure data to JSON"""
        data = {
            "fileName": self.filepath.name,
            "totalAtoms": self.get_total_atoms(),
            "chains": self.extract_chains(),
            "domains": self.extract_domains(),
            "ligands": self.extract_ligands()
        }
        
        return json.dumps(data, indent=2)


def main():
    """Main function"""
    if len(sys.argv) != 2:
        print(json.dumps({"error": "Usage: parse_protein.py <structure_file>"}))
        sys.exit(1)
        
    filepath = sys.argv[1]
    
    if not Path(filepath).exists():
        print(json.dumps({"error": f"File not found: {filepath}"}))
        sys.exit(1)
        
    try:
        parser = ProteinParser(filepath)
        parser.parse()
        print(parser.to_json())
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)


if __name__ == "__main__":
    main()